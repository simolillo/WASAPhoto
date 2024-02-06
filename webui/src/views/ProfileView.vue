<script>
// get user profile
export default {
	data: function() {
		return {
			username: "",
            photosCount: 0,
            followersCount: 0,
            followingCount: 0,
            isItMe: false,
            doIFollowUser: false,
            isInMyBannedList: false,
            amIBanned: false,

            userExists: false,
            userID: 0,
		}
	},
    watch: {
        // property to watch
        pathUsername(newUName, oldUName) {
            if (newUName !== oldUName){
                this.getUserProfile()
            }
        }
    },
    computed: {
        pathUsername(){
            return this.$route.params.username
        },
    },
    methods: {
        async getUserProfile() {
            if (this.$route.params.username === undefined){
                return
            }
            try {
                // profiles/username
				// GET /user/{Simo}
                let username = this.$route.params.username;
				let response = await this.$axios.get(`/user/${username}`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                this.userID = response.data;
                // GET /users/{1}/
                response = await this.$axios.get(`/users/${this.userID}/`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                let profile = response.data;
                console.log(profile)
                this.username = profile.username;
                this.photosCount = profile.photosCount;
                this.followersCount = profile.followersCount;
                this.followingCount = profile.followingCount;
                this.isItMe = profile.isItMe;
                this.doIFollowUser = profile.doIFollowUser;
                this.isInMyBannedList = profile.isInMyBannedList;
                this.amIBanned = profile.amIBanned;
                this.userExists = true;
            } catch (error) {
				const status = error.response.status;
        		const errorMessage = error.response.data;
        		alert(`Status (${status}): ${errorMessage}`);
            }
        },
		// async followBtn() {
        //     try {
        //         if (this.doIFollowUser) { 
        //             // DELETE /following/{1}
        //             await this.$axios.delete(`/following/${userID}`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
        //             this.followersCount -=1
        //         } else {
        //             // PUT /following/{1}
        //             await this.$axios.put(`/following/${userID}`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
        //             this.followersCount +=1
        //         }
        //         this.doIFollowUser = !this.followStatus
        //     } catch (error) {
        //         const status = error.response.status;
        // 		const errorMessage = error.response.data;
        // 		alert(`Status (${status}): ${errorMessage}`);
        //     }
		// },
        // async banBtn() {
        //     try {
        //         if (this.isInMyBannedList) {
        //             // DELETE /banned/{1}
        //             await this.$axios.delete(`/banned/${userID}`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
        //             this.getUserProfile();
        //         } else {
        //             // PUT /banned/{1}
        //             await this.$axios.put(`/banned/${userID}`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
        //             this.getUserProfile();
        //         }
        //     } catch (error) {
        //         const status = error.response.status;
        // 		const errorMessage = error.response.data;
        // 		alert(`Status (${status}): ${errorMessage}`);
        //     }
		// },
    },
    mounted() {
        this.getUserProfile();
    }
}
</script>

<template>

    <div class="container-fluid" v-if="userExists">
        <div class="row">
            <div class="col-12 d-flex justify-content-center">
                <div class="card w-50 container-fluid">
                    <div class="row">
                        <div class="col">
                            <div class="card-body d-flex justify-content-between align-items-center">
                                <h5 class="card-title p-0 me-auto mt-auto">@{{username}}</h5>

                                <button v-if="!isItMe && !isInMyBannedList" @click="followBtn" class="btn btn-success ms-2">
                                    {{doIFollowUser ? "Unfollow" : "Follow"}}
                                </button>

                                <button v-if="!isItMe" @click="banBtn" class="btn btn-danger ms-2">
                                    {{isInMyBannedList ? "Unban" : "Ban"}}
                                </button>
                            </div>
                        </div>
                    </div>

                    <div v-if="!isInMyBannedList" class="row mt-1 mb-1">
                        <div class="col-4 d-flex justify-content-start">
                            <h6 class="ms-3 p-0 ">Posts: {{photosCount}}</h6>
                        </div>
                    
                        <div class="col-4 d-flex justify-content-center">
                            <h6 class=" p-0 ">Followers: {{followersCount}}</h6>
                        </div>
                    
                        <div class="col-4 d-flex justify-content-end">
                            <h6 class=" p-0 me-3">Following: {{followingCount}}</h6>
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
                        <input id="fileUploader" type="file" class="profile-file-upload" @change="uploadFile" accept=".jpg, .png, .jpeg">
                        <label v-if="isItMe" class="btn my-btn-add-photo ms-2 d-flex align-items-center" for="fileUploader"> Add </label>
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
                    :owner="this.$route.params.id" 
                    :photo_id="photo.photo_id" 
                    :comments="photo.comments" 
                    :likes="photo.likes" 
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
