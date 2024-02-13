<script>
// get user profile
export default {
	data: function() {
		return {
            errormsg: null,

            // getUserProfile
			username: "",
            photosCount: 0,
            followersCount: 0,
            followingCount: 0,
            isItMe: false,
            doIFollowUser: false,
            isInMyBannedList: false,
            amIBanned: false,

            // getPhotosList
            photosList: [],

            // getFollowersList
            followersList: [],

            // getFollowingsList
            followingsList: [],

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
        pathUsername() {
            return this.$route.params.username
        },
    },
    methods: {
        async getUserProfile() {
            if (this.$route.params.username === undefined) {
                return
            }
            try {
                // /profiles/username
                // getUserId
                // GET /user/{username}
                let username = this.$route.params.username;
                let response = await this.$axios.get(`/user/${username}`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                this.userID = response.data;
                // GET /users/{uid}/
                response = await this.$axios.get(`/users/${this.userID}/`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                let profile = response.data;
                this.username = profile.username;
                this.photosCount = profile.photosCount;
                this.followersCount = profile.followersCount;
                this.followingCount = profile.followingCount;
                this.isItMe = profile.isItMe;
                this.doIFollowUser = profile.doIFollowUser;
                this.isInMyBannedList = profile.isInMyBannedList;
                this.amIBanned = profile.amIBanned;
                this.userExists = true;
                if (!this.isInMyBannedList && !this.amIBanned) {
                    await this.getPhotosList();
                    this.getFollowersList();
                    this.getFollowingsList();
                }
            } catch (error) {
                const status = error.response.status;
        		const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
            }
        },
		async followBtn() {
            try {
                if (this.doIFollowUser) { 
                     // DELETE /following/{uid}
                    await this.$axios.delete(`/following/${this.userID}`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                    this.getUserProfile();
                } else {
                    // PUT /following/{uid}
                    await this.$axios.put(`/following/${this.userID}`, null, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                    this.getUserProfile();
                }
                this.doIFollowUser = !this.doIFollowUser
            } catch (error) {
                const status = error.response.status;
        		const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
            }
		},
        async banBtn() {
            try {
                if (this.isInMyBannedList) {
                    // DELETE /banned/{uid}
                    await this.$axios.delete(`/banned/${this.userID}`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                    this.getUserProfile();
                } else {
                    // PUT /banned/{uid}
                    await this.$axios.put(`/banned/${this.userID}`, null, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                    this.getUserProfile();
                }
            } catch (error) {
                const status = error.response.status;
        		const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
            }
		},
        async uploadPhoto() {
            try {
                let file = document.getElementById('fileUploader').files[0];
                const reader = new FileReader();
                reader.readAsArrayBuffer(file); // stored in result attribute
                reader.onload = async () => {
                    // POST /photos/
                    let response = await this.$axios.post('/photos/', reader.result, {headers: {'Authorization': `${localStorage.getItem('token')}`, 'Content-Type': 'image/*'}});
                    this.photosList.unshift(response.data); // at the beginning of the list
                    this.photosCount += 1;
                }
            } catch (error) {
                const status = error.response.status;
        		const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
            }
        },
        async getPhotosList() {
            try {
                // GET /users/{uid}/photos/
                let response = await this.$axios.get(`/users/${this.userID}/photos/`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                this.photosList = response.data === null ? [] : response.data;
            } catch (error) {
				const status = error.response.status;
        		const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
            }
        },
        async getFollowersList() {
            try {
                // GET /users/{uid}/followers/
                let response = await this.$axios.get(`/users/${this.userID}/followers/`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                this.followersList = response.data === null ? [] : response.data;
            } catch (error) {
				const status = error.response.status;
        		const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
            }
        },
        async getFollowingsList() {
            try {
                // GET /users/{uid}/followings/
                let response = await this.$axios.get(`/users/${this.userID}/followings/`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                this.followingsList = response.data === null ? [] : response.data;
            } catch (error) {
				const status = error.response.status;
        		const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
            }
        },
        // on child event
        removePhotoFromList(photoID){
			this.photosList = this.photosList.filter(photo => photo.photoID != photoID);
            this.photosCount -= 1;
		},
        visitUser(username) {
            if (username != this.$route.params.username) {
                this.$router.push(`/profiles/${username}`);
            }
        }
    },
    mounted() {
        this.getUserProfile();
    }
}
</script>

<template>

    <UsersModal
    :modalID="'usersModalFollowers'" 
    :usersList="followersList"
    @visitUser="visitUser"
    />

    <UsersModal
    :modalID="'usersModalFollowing'" 
    :usersList="followingsList"
    @visitUser="visitUser"
    />

    <div class="container-fluid" v-if="userExists && !amIBanned">
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
                        <button class="col-4 d-flex justify-content-center btn-foll">
                            <h6 class="ms-3 p-0 ">Posts: {{photosCount}}</h6>
                        </button>
                    
                        <button class="col-4 d-flex justify-content-center btn-foll">
                            <h6 data-bs-toggle="modal" :data-bs-target="'#usersModalFollowers'">
                                Followers: {{followersList.length}}
                            </h6>
                        </button>
                    
                        <button class="col-4 d-flex justify-content-center btn-foll">
                            <h6 data-bs-toggle="modal" :data-bs-target="'#usersModalFollowing'">
                                Following: {{followingsList.length}}
                            </h6>
                        </button>
                    </div>
                </div>
            </div>
        </div>


        <div class="row">
            <div class="container-fluid mt-3">
                <div class="row ">
                    <div class="col-12 d-flex justify-content-center">
                        <h2>Posts</h2>
                        <input id="fileUploader" type="file" class="profile-file-upload" @change="uploadPhoto" accept=".jpg, .png">
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
                <div v-if="!isInMyBannedList && photosCount>0">
                    <Photo v-for="photo in photosList"
                    :key="photo.photoID"
                    :photoID="photo.photoID"
                    :authorID="photo.authorID"
                    :authorUsername="this.username"
                    :date="photo.date"
                    :likesListParent="photo.likesList"
                    :commentsListParent="photo.commentsList"
                    :isItMe="isItMe"
                    @removePhoto="removePhotoFromList"
                    />
                </div>
                
                <div v-if="!isInMyBannedList && photosCount==0" class="mt-5 ">
                    <h2 class="d-flex justify-content-center" style="color: white;">No posts yet</h2>
                </div>
            </div>
        </div>
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    
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
.btn-foll{
    background-color: transparent;
    border: none;
    padding: 5px;
}
</style>
