<script>
export default {
	data: function() {
		return {
			errormsg: null,

			userExists: false,
			banStatus: false,

            nickname: "",


			followStatus: false,
			currentIsBanned: false,

			followerCnt: 0,
			followingCnt:0,
			postCnt:0,

			photos: [],
            following: [],
            followers: [],
		}
	},
    // Monitora i cambiamenti nell'ID dell'utente nel percorso e ricarica le informazioni del profilo se cambia.
    watch:{
        currentPath(newid,oldid){
            if (newid !== oldid){
                this.loadInfo()
            }
        },
    },

	computed:{
        // Ottiene l'ID dell'utente dal percorso corrente.
        currentPath(){
            return this.$route.params.uid
        },
        
        // Controlla se l'utente visualizzato è lo stesso che ha effettuato l'accesso.
		sameUser(){
			return this.$route.params.uid === sessionStorage.getItem('token')
		},
	},

	methods: {

        async uploadFile(){
            let fileInput = document.getElementById('fileUploader')

            const file = fileInput.files[0];
            const reader = new FileReader();

            reader.readAsArrayBuffer(file);

            reader.onload = async () => {
                // Post photo: /users/:id/photos
                let response = await this.$axios.post("/photos/", reader.result, {
                    headers: {
                    'Content-Type': file.type
                    },
                })
                //console.log(response)
                /*
                this.photos.unshift({
                    owner: response.data.owner,
                    date: response.data.date,
                    photo_id: response.data.photo_id,
                    likes: response.data.likes,
                    comments: response.data.comments,
                })
                */
                this.photos.unshift(response.data)
                this.postCnt += 1
            location.reload();
            };
        },

		async followClick(){
            try{
                if (this.followStatus){ 
                    await this.$axios.delete("/following/"+this.$route.params.uid);
                    this.followerCnt -=1
                }else{
                    await this.$axios.put("/following/"+this.$route.params.uid);
                    this.followerCnt +=1
                }
                this.followStatus = !this.followStatus
            }catch (e){
                this.errormsg = e.toString();
            }
            
		},

		async banClick(){
            try{
                if (this.banStatus){
                    await this.$axios.delete("/banned/"+this.$route.params.uid);
                    this.loadInfo()
                }else{
                    await this.$axios.put("/banned/"+this.$route.params.uid);
                    this.followStatus = false
                }
                this.banStatus = !this.banStatus
            }catch(e){
                this.errormsg = e.toString();
            }
		},
        async followersList(uid) {
			try {
				this.errormsg = null
				let response = await this.$axios.get("/users/" + uid + "/followers/")

				if (response.data != null){
					return response.data
				}

				return [];
			} catch (e) {
				this.errormsg = e.toString()
				return [];
			}
		},
        async followingsList(uid) {
			try {
				this.errormsg = null
				let response = await this.$axios.get("/users/" + uid + "/followings/")

				if (response.data != null){
					return response.data
				}

				return [];
			} catch (e) {
				this.errormsg = e.toString()
				return [];
			}
		},
        async photosList(uid) {
			try {
				this.errormsg = null
				let response = await this.$axios.get("/users/" + uid + "/photos/")

				if (response.data != null){
					return response.data
				}

				return [];
			} catch (e) {
				this.errormsg = e.toString()
				return [];
			}
		},
		async loadInfo(){
            if (this.$route.params.uid === undefined){
                return
            }

			try{
                // Get user profile: /users/:id
				let response = await this.$axios.get("/users/"+this.$route.params.uid);

                this.banStatus = false
                this.userExists = true
                this.currentIsBanned = false

                if (response.status === 206){
                    this.banStatus = true
                    return
                }

				if (response.status === 204){
					this.userExists = false
				}
				
                this.nickname = response.data.username
				this.followerCnt = response.data.FollowersCount
				this.followingCnt = response.data.FollowingCount
				this.postCnt = response.data.PhotosCount
				this.followStatus = response.data.IsFollowedByViewer
                this.photos = this.photosList(this.$route.params.uid)
                this.followers = this.followersList(this.$route.params.uid)
                this.following = this.followingsList(this.$route.params.uid)

			}catch(e){
				this.currentIsBanned = true
			}
		},

        goToSettings(){
            this.$router.push('/settings')
        },

        removePhotoFromList(photo_id){
			this.photos = this.photos.filter(item => item.photo_id !== photo_id)
		},
        async commentsList(pid) {
			try {
				this.errormsg = null
				let response = await this.$axios.get("/photos/" + pid + "/comments/")

				if (response.data != null){
					return response.data
				}

				return [];
			} catch (e) {
				this.errormsg = e.toString()
				return [];
			}
		},
		async likesList(pid) {
			try {
				this.errormsg = null
				let response = await this.$axios.get("/photos/" + pid + "/likes/")

				if (response.data != null){
					return response.data
				}

				return [];
			} catch (e) {
				this.errormsg = e.toString()
				return [];
			}
		}
	},
    // Carica le informazioni del profilo quando il componente viene montato.
	async mounted(){
		await this.loadInfo()
	},

}
</script>

<template>

    <div class="container-fluid" v-if="!currentIsBanned && userExists">

        <div class="row">
            <div class="col-12 d-flex justify-content-center">
                <div class="card w-50 container-fluid">

                    <div class="row">
                        <div class="col">
                            <div class="card-body d-flex justify-content-between align-items-center">
                                <h5 class="card-title p-0 me-auto mt-auto">{{nickname}} @{{this.$route.params.uid}}</h5>

                                <button v-if="!sameUser && !banStatus" @click="followClick" class="btn btn-success ms-2">
                                    {{followStatus ? "Unfollow" : "Follow"}}
                                </button>

                                <button v-if="!sameUser" @click="banClick" class="btn btn-danger ms-2">
                                    {{banStatus ? "Unban" : "Ban"}}
                                </button>

                                <button v-else class="my-trnsp-btn ms-2" @click="goToSettings">
                                    <!--Settings  <font-awesome-icon icon="fa-solid fa-gear" />-->
                                    <i class="my-nav-icon-gear fa-solid fa-gear"></i>
                                </button>
                            </div>
                        </div>
                    </div>

                    <div v-if="!banStatus" class="row mt-1 mb-1">
                        <div class="col-4 d-flex justify-content-start">
                            <h6 class="ms-3 p-0 ">Posts: {{postCnt}}</h6>
                        </div>
                    
                        <div class="col-4 d-flex justify-content-center">
                            <h6 class=" p-0 ">Followers: {{followerCnt}}</h6>
                        </div>
                    
                        <div class="col-4 d-flex justify-content-end">
                            <h6 class=" p-0 me-3">Following: {{followingCnt}}</h6>
                        </div>
                    </div>
                </div>
            </div>
        </div>


        <div class="row">

            <div class="container-fluid mt-3">

                <div class="row ">
                    <div class="col-12 d-flex justify-content-center">
                        <h2>Posts</h2>
                        <input id="fileUploader" type="file" class="profile-file-upload" @change="uploadFile" accept=".jpg, .png">
                        <label v-if="sameUser" class="btn my-btn-add-photo ms-2 d-flex align-items-center" for="fileUploader"> Add </label>
                    </div>
                </div>

                <div class="row ">
                    <div class="col-3"></div>
                    <div class="col-6">
                        <hr class="border border-dark">
                    </div>
                    <div class="col-3"></div>
                </div>
            </div>
        </div>

        <div class="row">
            <div class="col">

                <div v-if="!banStatus && postCnt>0">
                    <Photo v-for="(photo,index) in photos" 
                    :key="index" 
                    :owner="photo.authorID" 
                    :photo_id="photo.photoID" 
                    :comments="commentsList(photo.photoID)" 
                    :likes="likesList(photo.photoID)" 
                    :upload_date="photo.date" 
                    :isOwner="sameUser" 
                    
                    @removePhoto="removePhotoFromList"
                    />

                </div>
                
                <div v-else class="mt-5 ">
                    <h2 class="d-flex justify-content-center" style="color: white;">No posts yet</h2>
                </div>

            </div>
        </div>

    
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    </div>
    <div v-else class="h-25 ">
        <PageNotFound />
    </div>
    

</template>

<style>
.profile-file-upload{
    display: none;
}

.my-nav-icon-gear{
    color: grey;
}
.my-nav-icon-gear:hover{
    transform: scale(1.3);
}

.my-btn-add-photo{
    background-color: green;
    border-color: grey;
}
.my-btn-add-photo:hover{
    color: white;
    background-color: green;
    border-color: grey;
}

</style>
