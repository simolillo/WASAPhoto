<script>
// stream
export default {
	data: function () {
		return {
			photos: [],
		}
	},
	methods: {
		async getMyStream() {
			try {
				// GET /stream
                let response = await this.$axios.get('/stream', {headers: {'Authorization': `${localStorage.getItem('token')}`}});
				this.photos = response.data === null ? [] : response.data;
				console.log(this.photos)
			} catch (error) {
				const status = error.response.status;
        		const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
            }
		},
		async getAuthorUsername(authorID) {
			try {
				// GET /stream
                let response = await this.$axios.get('/stream', {headers: {'Authorization': `${localStorage.getItem('token')}`}});
				this.photos = response.data === null ? [] : response.data;
				console.log(this.photos)
			} catch (error) {
				const status = error.response.status;
        		const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
            }
		}
	},
	async mounted() {
		await this.getMyStream()
	}
}
</script>

<template>
	<div class="container-fluid">
		<div class="row">
			<Photo v-for="photo in photos"
			:key="photo.photoID"
			:photoID="photo.photoID"
			:authorID="photo.authorID"
			:authorUsername="photo.authorUsername"
			:date="photo.date"
			:likesListParent="photo.likesList"
			:commentsListParent="photo.commentsList"
			:isItMe="false"
			/>
		</div>
		<div v-if="photos.length === 0" class="row">
			<h1 class="d-flex justify-content-center mt-5" style="color: rgb(0, 0, 0);">There's no content yet, follow somebody!</h1>
		</div>
	</div>
</template>

<style>
</style>
