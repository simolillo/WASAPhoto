<script>
export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			banned: [],
			userId: null,
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

				let banResponse = await this.$axios.get(
					`/users/${this.userId}/banned/`,
					{
						headers: {
							Authorization:
								"Bearer " + localStorage.getItem("token"),
						},
					}
				);
				this.banned = banResponse.data;
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
		<h1>Ban list</h1>

		<ErrorMsg
			class="error-container"
			v-if="errormsg"
			:msg="errormsg"
		></ErrorMsg>

		<LoadingSpinner v-if="loading"></LoadingSpinner>

		<div v-if="banned?.length === 0">
			<div class="card">
				<div class="card-body">
					<p>No banned in the database.</p>
				</div>
			</div>
		</div>

		<div v-if="!loading">
			<div v-for="user in banned" :key="user.id">
				<a
					href="javascript:"
					class="text-muted mb-1 larger-text"
					style="text-decoration: none"
					@click="gotoAccount(this.userId, user.bannedUser)"
					>{{ user.bannedUser }}</a
				>
			</div>
		</div>
	</div>
</template>
