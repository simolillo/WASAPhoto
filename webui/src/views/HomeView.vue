<script>
export default {
	data: function () {
		return {
			errormsg: null,
			photos: [],
			likes: [],
		}
	},
	methods: {
		async loadStream() {
			try {
				this.errormsg = null
				let response = await this.$axios.get("/stream")

				if (response.data != null){
					this.photos = response.data
				}
				
			} catch (e) {
				this.errormsg = e.toString()
			}
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

	async mounted() {
		await this.loadStream()
	}

}
</script>

<template>
	<div class="container-fluid">
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<div class="row">
			<Photo
				v-for="(photo,index) in photos"
				:key="index"
				:owner="photo.authorID"
				:photo_id="photo.photoID"
				:comments="commentsList(photo.photoID)"
				:likes="likesList(photo.photoID)"
				:upload_date="photo.date"
			/>
		</div>

		<div v-if="photos.length === 0" class="row ">
			<h1 class="d-flex justify-content-center mt-5" style="color: white;">There's no content yet, follow somebody!</h1>
		</div>
	</div>
</template>

<style>
</style>
