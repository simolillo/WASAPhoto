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
		}
	},
    methods: {
        async getUserProfile() {
            try {
                // profiles/username
				// GET /user/{Simo}
                let username = this.$route.params.username;
				let response = await this.$axios.get(`/user/${username}`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                let userID = response.data;
                // GET /users/{1}/
                response = await this.$axios.get(`/users/${userID}/`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
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
            } catch (error) {
                alert('hola amigo')
				const status = error.response.status;
        		const errorMessage = error.response.data;
        		//alert(`Status (${status}): ${errorMessage}`);
            }
        }
    },
    mounted() {
        this.getUserProfile();
    }
}
</script>

<template>
    <h1 v-if="isItMe">
        username: {{ username }}
        photosCount: {{ photosCount }}
        followersCount: {{ followersCount }}
        followingCount: {{ followingCount }}
        isItMe: {{ isItMe }}
        doIFollowUser: {{ doIFollowUser }}
        isInMyBannedList: {{ isInMyBannedList }}
        amIBanned: {{ amIBanned }}
    </h1>
</template>

<style>
</style>
