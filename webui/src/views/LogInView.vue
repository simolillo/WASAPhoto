<script>
export default {
	data: function() {
		return {
			errormsg: null,
			user: {
                ID: 0,
                Name: "",                
            }
		}
	},
	methods: {
		async login() {
			this.errormsg = null;
			try {
				let response = await this.$axios.post("/session",{username: this.user.Name});
                this.user = response.data;
                sessionStorage.setItem('token', this.user.ID)
				this.$router.replace("/home")
				this.$emit('updatedLoggedChild',true)
			} catch (e) {
				this.errormsg = e.toString();
			}
		},
	},
	mounted(){
		if (sessionStorage.getItem('token')){
			this.$router.replace("/home")
		}
	},
	
}
</script>

<template>
	<div class="container-fluid h-100 m-0 p-0 login">

		<div class="row ">
			<div class="col">
				<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
			</div>
		</div>

		<div class="row h-100 w-100 m-0">
			
			<form @submit.prevent="login" class="d-flex flex-column align-items-center justify-content-center p-0">

				<div class="row mt-2 mb-3 border-bottom">
					<div class="col">
						<h2 class="login-title">WASAPhoto Login</h2>
					</div>
				</div>

				<div class="row mt-2 mb-3">
					<div class="col">
						<input 
						type="text" 
						class="form-control" 
						v-model="username" 
						maxlength="16"
						minlength="3"
						placeholder="Your identifier" />
					</div>
				</div>

				<div class="row mt-2 mb-5 ">
					<div class="col ">
						<button class="btn btn-dark" :disabled="username == null || username.length >16 || username.length <3 || username.trim().length<3"> 
						Register/Login 
						</button>
					</div>
				</div>
			</form>
		</div>
	</div>
</template>

<style>
.login {
    background-image: url("../assets/images/BackgroundLogin.jpg");
    height: 100vh;
}

.login-title {
    color: black;
}
</style>
