<script>
export default {
	data() {
		return {
            errormsg: null,

			photoURL: "",
			liked: false,
		}
	},

	props: ['photoID','authorID','authorUsername','date','likesList','commentsList','isItMe'], 

	methods: {
		getPhoto() {
			// GET /photos/{1}/
			this.photoURL = __API_URL__ + `/photos/${this.photoID}/`;
		},
		async deletePhoto() {
			try {
				// DELETE /photos/{1}/
				await this.$axios.delete(`/photos/${this.photoID}/`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
				this.$emit('removePhoto', this.photoID); // signal to parent
			} catch (error) {
                const status = error.response.status;
        		const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
            }
		},
		visitAuthorProfile() {
            // /profiles/:username
			this.$router.push(`/profiles/${this.authorUsername}`)
		},
		async likeToggle() {
			try {
				if (!this.liked){
					// Put like: /users/:id/photos/:photo_id/likes/:like_id"
					await this.$axios.put("/users/"+ this.owner +"/photos/"+this.photo_id+"/likes/"+ bearer)
					this.allLikes.push({
						user_id: bearer,
						nickname: bearer
					})

				}else{
					// Delete like: /users/:id/photos/:photo_id/likes/:like_id"
					await this.$axios.delete("/users/"+ this.owner  +"/photos/"+this.photo_id+"/likes/"+ bearer)
					this.allLikes.pop()
				}

				this.liked = !this.liked;
			}catch(e){
				//
			}
      		
    	},

		removeCommentFromList(value){
			this.allComments = this.allComments.filter(item=> item.comment_id !== value)
		},

		addCommentToList(comment){
			this.allComments.push(comment)
		},
	},
	
	async mounted(){
		await this.getPhoto()

		if (this.likes != null){
			this.allLikes = this.likes
		}

		if (this.likes != null){
			this.liked = this.allLikes.some(obj => obj.user_id === sessionStorage.getItem('token'))
		}
		if (this.comments != null){
			this.allComments = this.comments
		}
		
		
	},

}
</script>

<template>
	<div class="container-fluid mt-3 mb-5 ">

        <LikeModal :modal_id="'like_modal'+photo_id" 
		:likes="allLikes" />

        <CommentModal :modal_id="'comment_modal'+photo_id" 
		:comments_list="allComments" 
		:photo_owner="owner" 
		:photo_id="photo_id"

		@eliminateComment="removeCommentFromList"
		@addComment="addCommentToList"
		/>

        <div class="d-flex flex-row justify-content-center">

            <div class="card my-card">
                <div class="d-flex justify-content-end">

                    <button v-if="isOwner" class="my-trnsp-btn my-dlt-btn me-2" @click="deletePhoto">
						<!--Delete-->
						<i class="fa-solid fa-trash w-100 h-100"></i>
					</button>

                </div>
                <div class="d-flex justify-content-center photo-background-color">
                    <img :src="photoURL" class="card-img-top img-fluid">
                </div>

                <div class="card-body">

                    <div class="container">

                        <div class="d-flex flex-row justify-content-end align-items-center mb-2">

							<button class="my-trnsp-btn m-0 p-1 me-auto" @click="visitAuthorProfile">
                            	<i> From {{owner}}</i>
							</button>

                            <button class="my-trnsp-btn m-0 p-1 d-flex justify-content-center align-items-center">
                                <i @click="likeToggle" :class="'me-1 my-heart-color w-100 h-100 fa '+(liked ? 'fa-heart' : 'fa-heart-o') "></i>
                                <i data-bs-toggle="modal" :data-bs-target="'#like_modal'+photo_id" class="my-comment-color ">
                                    {{allLikes.length}}
                                </i>
                            </button>

                            <button class="my-trnsp-btn m-0 p-1  d-flex justify-content-center align-items-center" 
							data-bs-toggle="modal" :data-bs-target="'#comment_modal'+photo_id">

                                <i class="my-comment-color fa-regular fa-comment me-1" @click="commentClick"></i>
                                <i class="my-comment-color-2"> {{allComments != null ? allComments.length : 0}}</i>

                            </button>
                        </div>

                        <div class="d-flex flex-row justify-content-start align-items-center ">
                            <p> Uploaded on {{upload_date}}</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style>
.photo-background-color{
	background-color: grey;
}

.my-card{
	width: 27rem;
	border-color: black;
	border-width: thin;
}

.my-heart-color{
	color: grey;
}
.my-heart-color:hover{
	color: red;
}

.my-comment-color {
	color: grey;
}
.my-comment-color:hover{
	color: black;
}

.my-comment-color-2{
	color:grey
}

.my-dlt-btn{
	font-size: 19px;
}
.my-dlt-btn:hover{
	font-size: 19px;
	color: var(--color-red-danger);
}
</style>
