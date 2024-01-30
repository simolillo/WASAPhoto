<script>
import "../styles.css";

export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			photo: null,
			userId: localStorage.getItem("token"),
			username: localStorage.getItem("username"),
		};
	},
	methods: {
		load() {
			return load;
		},
		async refresh() {
			this.loading = true;
			this.errormsg = null;
		},
		uploadPhoto: async function () {
			this.loading = true;
			this.errormsg = null;
			try {
				let fileInput = document.getElementById("photo");
				let file = fileInput.files[0];
				let path = file.name;
				await this.$axios.post(
					`/users/${this.userId}/photos/`,
					{
						path: path,
					},
					{
						headers: {
							Authorization:
								"Bearer " + localStorage.getItem("token"),
						},
					}
				);
				// console.log(this.photo.path);
				this.$router.push({
					name: "MyAccount",
					params: { username: this.username },
				});
			} catch (e) {
				this.errormsg =
					"Error uploading photo: " + e.response.data.message;
			}
			this.loading = false;
		},
	},

	mounted() {
		this.refresh();
	},
};
</script>

<template>
	<div class="background">
		<div class="container py-5">
			<div class="row justify-content-center align-items-center">
				<div class="col-md-6">
					<div class="card bg-white text-dark rounded-3">
						<div class="card-body p-5 text-center">
							<h2 class="fw-bold mb-4 text-uppercase">
								Upload Photo
							</h2>
							<p class="text-muted">Please select your photo.</p>
							<div class="form-group">
								<input
									type="file"
									id="photo"
									class="form-control form-control-lg rounded-pill mb-3"
									accept="image/*"
									required
								/>
							</div>
							<div class="d-grid gap-3">
								<button
									class="btn btn-primary rounded-pill"
									type="submit"
									@click="uploadPhoto"
									style="background-color: #2e4a78"
								>
									Upload
								</button>
								<LoadingSpinner v-if="loading" />
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
		<ErrorMsg class="error-container" v-if="errormsg" :msg="errormsg" />
	</div>
</template>

<style></style>
