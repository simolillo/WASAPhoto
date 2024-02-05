<script>
// update username
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
		}
	},
    methods: {
        async getUserProfile() {
            try {
				// PUT /settings
				let username = document.getElementById('username').value;
				if (!username.match("^[a-zA-Z][a-zA-Z0-9_]{2,15}$")) {
                alert("Invalid username: 3 - 16 characters; first character must be a letter; only letters, numbers and underscores allowed");
                return;
				}
				let response = await this.$axios.put('/settings', {username: username}, {headers: {'Authorization': `${localStorage.getItem('token')}`, 'Content-Type': 'application/json'}});
				let user = response.data // userID, username
				localStorage.setItem('token', user.userID);
				localStorage.setItem('username', user.username);
                alert(`Username correctly updated to: ${user.username}`)
            } catch (error) {
				const status = error.response.status;
        		const errorMessage = error.response.data;
        		alert(`Status (${status}): ${errorMessage}`);
            }
        }
    },
    mounted() {
        this.finishedMounted = true;
    }
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
                                <h5 class="card-title p-0 me-auto mt-auto">{{nickname}} @{{this.$route.params.id}}</h5>

                                <button v-if="!sameUser && !banStatus" @click="followClick" class="btn btn-success ms-2">
                                    {{followStatus ? "Unfollow" : "Follow"}}
                                </button>

                                <button v-if="!sameUser" @click="banClick" class="btn btn-danger ms-2">
                                    {{banStatus ? "Unban" : "Ban"}}
                                </button>

                                <button v-else class="my-trnsp-btn ms-2" @click="goToSettings">
                     
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
</style>
