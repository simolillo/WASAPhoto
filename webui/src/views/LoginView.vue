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

                let response = await this.$axios.post("/session", {"username": username}, {headers: {'Content-Type': 'application/json'}});
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
		if (localStorage.getItem("token")){
			this.$router.replace("/stream")
		}
	}
}
</script>

<template>
</template>

<style></style>
