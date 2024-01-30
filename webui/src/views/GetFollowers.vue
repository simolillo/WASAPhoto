<script>
export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			followers: [],
			userId: null,
			username: null,
		};
	},
	methods: {
		load() {
			return load;
		},
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				this.userId = this.$route.params.userId;
				this.username = this.$route.params.username;

				let apiUrl = `/profile/${this.username}/`;
				let response = await this.$axios.get(apiUrl, {
					headers: {
						Authorization:
							"Bearer " + localStorage.getItem("token"),
					},
				});
				this.followers = response.data.followers;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async gotoAccount(id, username) {
			this.$router.push({
				name: "MyAccount",
				params: { username: username },
			});
		},
	},

	mounted() {
		this.refresh();
	},
};
</script>

<template>
	<div>
		<h1>Followers list</h1>

		<ErrorMsg
			class="error-container"
			v-if="errormsg"
			:msg="errormsg"
		></ErrorMsg>

		<LoadingSpinner v-if="loading"></LoadingSpinner>
		<div v-if="followers?.length === 0">
			<div class="card">
				<div class="card-body">
					<p>No followers in the database.</p>
				</div>
			</div>
		</div>
		<div v-if="!loading">
			<div v-for="f in followers" :key="f">
				<a
					href="javascript:"
					class="text-muted mb-1 larger-text"
					style="text-decoration: none"
					@click="gotoAccount(this.userId, f)"
					>{{ f }}</a
				>
			</div>
		</div>
	</div>
</template>
