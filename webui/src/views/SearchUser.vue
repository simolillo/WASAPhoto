<script>
export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			otheruser: null,
			userId: localStorage.getItem("token"),
		};
	},
	methods: {
		async findUser() {
			this.loading = true;
			this.errormsg = null;
			try {
				let url = `/profile/${this.otheruser}/`;
				let response = await this.$axios.get(url, {
					headers: {
						Authorization:
							"Bearer " + localStorage.getItem("token"),
					},
				});
				this.$router.push({
					name: "MyAccount",
					params: { username: this.otheruser },
				});
			} catch (e) {
				if (e.response.status == 401) {
					this.errormsg = "You can not acces " + this.otheruser;
				} else if (e.response.status == 404) {
					// this.$router.push({ name: "Login" });
					this.errormsg = "User not found " + this.otheruser;
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
								Search User
							</h2>
							<p class="text-muted">
								Please enter the username of the user that you
								are searching for.
							</p>
							<div class="form-group">
								<input
									type="text"
									id="otheruser"
									class="form-control form-control-lg rounded-pill mb-3"
									v-model="otheruser"
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
									@click="findUser"
									style="background-color: #2e4a78"
								>
									Search User
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
