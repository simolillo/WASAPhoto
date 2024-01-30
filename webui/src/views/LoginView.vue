<script>
export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			username: null,
			user: {
				userId: 0,
				username: "",
			},
		};
	},
	methods: {
		createUser: async function () {
			this.loading = true;
			this.errormsg = null;
			try {
				const response = await this.$axios.post("/session", {username: this.username,});
				this.user = response.data;
				localStorage.setItem("token", this.user.id);
				localStorage.setItem("username", this.user.username);
				await this.$router.push({
					name: "MyAccount",
					params: { username: this.username },
				});
			} catch (e) {
				this.errormsg = e.toString();
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
							<h2 class="fw-bold mb-4 text-uppercase">Welcome</h2>
							<p class="text-muted">
								Please enter your username.
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
							<div class="d-grid gap-3" @click="createUser">
								<button
									v-if="!loading"
									class="btn btn-primary rounded-pill"
									type="submit"
									style="background-color: #2e4a78"
								>
									Login
								</button>
								<LoadingSpinner v-if="loading" />
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
		<ErrorMsg v-if="errormsg" :msg="errormsg" />
	</div>
</template>

<style></style>
