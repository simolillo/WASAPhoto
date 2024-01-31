<script>
export default {
data() {
    return {
    username: "",
    errormsg: null,
    loading: false,
    };
},
methods: {
    async login() {
    this.loading = true;
    this.errormsg = null;

    try {
        // Assuming you have a doLogin method in your API
        let response = await this.$axios.post("/session", {
        username: this.username,
        });

        // Assuming the response contains user data
        let userData = response.data.user;

        // Do something with the user data (e.g., store in Vuex)
        console.log("User ID:", userData.userID);
        console.log("Username:", userData.username);

        // Redirect to home page or perform other actions as needed
        // Example: this.$router.push("/home");
    } catch (e) {
        this.errormsg = e.toString();
    }

    this.loading = false;
    },
},
};
</script>

<template>
    <div>
      <h1 class="h2">Login</h1>
  
      <form @submit.prevent="login">
        <div class="mb-3">
          <label for="username" class="form-label">Username</label>
          <input type="text" class="form-control" v-model="username" required>
        </div>
  
        <button type="submit" class="btn btn-primary">Login</button>
      </form>
  
      <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    </div>
  </template>
  
  
  <style>
  /* Add any styling if needed */
  </style>
  