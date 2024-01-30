<script>
export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			username: null,
			userId: localStorage.getItem("token"),
		};
	},
	methods: {
		async updateUsername() {
			this.loading = true;
			this.errormsg = null;
			try {
				const response = await this.$axios.put(
					`/users/${this.userId}/`,
					{
						username: this.username,
					},
					{
						headers: {
							Authorization:
								"Bearer " + localStorage.getItem("token"),
						},
					}
				);
				const responseData = response.data;
				localStorage.setItem("username", this.username);
				this.$router.push({
					name: "MyAccount",
					params: { username: this.username },
				});
			} catch (e) {
				if (e.response.status == 404) {
					this.errormsg = "User not found " + this.username;
				} else if (e.response.status == 409) {
					this.errormsg =
						"This username already exists " + this.username;
				} else if (e.response.status == 400) {
					this.errormsg = "Wrong format received " + this.username;
				} else {
					this.errormsg = e.toString();
				}
			}
			this.loading = false;
		},
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
								Update my username
							</h2>
							<p class="text-muted">
								Please enter you new username.
							</p>
							<div class="form-group">
								<input
									type="text"
									id="username"
									class="form-control form-control-lg rounded-pill mb-3"
									v-model="username"
									placeholder="Username"
									required
									minlength="3"
									maxlength="16"
								/>
							</div>
							<div class="d-grid gap-3">
								<button
									v-if="!loading"
									class="btn btn-primary rounded-pill"
									type="submit"
									@click="updateUsername"
									style="background-color: #2e4a78"
								>
									Update User
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
