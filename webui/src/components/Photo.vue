<script>
export default {
	data() {
		return {
            errormsg: null,

			photoURL: "",
			liked: false,
            likesList: [],
            commentsList: [],
		}
	},

	props: ['photoID','authorID','authorUsername','date','likesListParent','commentsListParent','isItMe'], 

	methods: {
		getPhoto() {
			// GET /photos/{pid}/
			this.photoURL = __API_URL__ + `/photos/${this.photoID}/`;
		},
		async deletePhoto() {
			try {
				// DELETE /photos/{pid}/
				await this.$axios.delete(`/photos/${this.photoID}/`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
				this.$emit("removePhoto", this.photoID);
			} catch (error) {
                const status = error.response.status;
        		const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
                alert(this.errormsg);
            }
		},
		visitAuthorProfile() {
            // /profiles/:username
			this.$router.push(`/profiles/${this.authorUsername}`);
		},
		async likeToggle() {
			try {
				if (!this.liked) {
					// PUT /likes/{pid}
                    await this.$axios.put(`/likes/${this.photoID}`, null, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
					this.likesList.push({userID: localStorage.getItem('token'), username: localStorage.getItem('username')});
				} else {
					// DELETE /likes/{pid}
                    await this.$axios.delete(`/likes/${this.photoID}`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                    this.likesList = this.likesList.filter(user => user.userID != localStorage.getItem('token'));
				}
				this.liked = !this.liked;
			} catch (error) {
                const status = error.response.status;
        		const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
                alert(this.errormsg);
            }
    	},
        // on child event
		removeCommentFromList(commentID) {
			this.commentsList = this.commentsList.filter(comment => comment.commentID != commentID);
		},
		addCommentToList(comment){
			this.commentsList.unshift(comment); // at the beginning of the list
		},
        visitLiker(username) {
            if (username != this.$route.params.username) {
                document.querySelector('.modal-backdrop').remove();
                document.querySelector('.modal').remove();
                document.body.style.overflow = 'auto';
                this.$router.push(`/profiles/${username}`);
            }
        }
	},
	async mounted() {
        this.getPhoto()
        // it is a promise
        if (this.likesListParent != null) {
            this.likesList = this.likesListParent
        }
        if (this.commentsListParent != null) {
            this.commentsList = this.commentsListParent
        }
		this.liked = this.likesList.some(user => user.userID == localStorage.getItem('token'));
	},
}
</script>

<template>
	<div class="container-fluid mt-3 mb-5 ">

        <UsersModal
        :modalID="'likesModal'+photoID" 
		:usersList="likesList"
        @visitUser="visitLiker"
        />

        <CommentModal
        :modalID="'commentModal'+photoID" 
		:commentsList="commentsList" 
		:isItMe="isItMe" 
		:photoID="photoID"
		@removeComment="removeCommentFromList"
		@addComment="addCommentToList"
		/>

        <div class="d-flex flex-row justify-content-center">
            <div class="card my-card">
                <div class="d-flex justify-content-end">
                    <button v-if="isItMe" class="my-trnsp-btn my-dlt-btn me-2" @click="deletePhoto">
						<!--trash bin-->
						<i class="fa-solid fa-trash w-100 h-100"></i>
					</button>
                </div>
                <!--photo-->
                <div class="d-flex justify-content-center photo-background-color">
                    <img :src="photoURL" class="card-img-top img-fluid">
                </div>
                <div class="card-body">
                    <div class="container">
                        <div class="d-flex flex-row justify-content-end align-items-center mb-2">
                            <!--author-->
							<button class="my-trnsp-btn m-0 p-1 me-auto" @click="visitAuthorProfile">
                            	<i> From {{authorUsername}}</i>
							</button>
                            <!--like-->
                            <button class="my-trnsp-btn m-0 p-1 d-flex justify-content-center align-items-center">
                                <i @click="likeToggle" :class="'me-1 my-heart-color w-100 h-100 fa '+(liked ? 'fa-heart' : 'fa-heart-o')"></i>
                                <i data-bs-toggle="modal" :data-bs-target="'#likesModal'+photoID" class="my-comment-color ">
                                    {{likesList.length}}
                                </i>
                            </button>
                            <!--comment-->
                            <button class="my-trnsp-btn m-0 p-1  d-flex justify-content-center align-items-center" 
							data-bs-toggle="modal" :data-bs-target="'#commentModal'+photoID">
                                <i class="my-comment-color fa-regular fa-comment me-1"></i>
                                <i class="my-comment-color-2"> {{commentsList.length}}</i>
                            </button>
                        </div>
                        <div class="d-flex flex-row justify-content-start align-items-center ">
                            <p> Uploaded on {{date}}</p>
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
