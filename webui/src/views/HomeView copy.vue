<script>
export default {
	data: function() {
		return {
            errormsg: null,
			user: {
                ID: 0,
                Name: "",
            },
		}
	},
	methods: {
		async login() {
            this.errormsg = null;
            try {
                let username = document.getElementById("login-form").value;
                if (!username.match("^[a-zA-Z][a-zA-Z0-9_]{2,32}$")) {
                    alert("Invalid username. First character must be a letter, only letters, numbers and underscores allowed.");
                    return;
                }

                let response = await this.$axios.post("/session", {
                    "username": this.user.Name
                });
                this.user = response.data;

                localStorage.setItem("token", this.user.ID);
                localStorage.setItem("username", this.user.Name);
                this.$router.replace("/stream");
            } catch (e) {
                this.errormsg = e.toString();
            }
        },
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
    <div class="container text-center pt-3 pb-2 border-bottom">
        <h2>Login</h2>
    </div>

    <div class="h-75 d-flex align-items-center justify-content-center">
        <form class="border border-dark p-5 rounded shadow-lg" @submit.prevent="login">
            <!-- Username input -->
            <div class="form-outline mb-4">
                <input type="text" id="login-form" class="form-control" pattern="^[a-zA-Z][a-zA-Z0-9_]{2,32}$" />
                <label class="form-label" for="login-form">Username</label>
            </div>
            <!-- Submit button -->
            <button type="submit" class="btn btn-primary btn-block mb-4">Sign in</button>
        </form>
    </div>
</template>

<style></style>
