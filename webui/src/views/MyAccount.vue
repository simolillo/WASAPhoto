<script>
import "../styles.css";

export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			user: null,
			userId: localStorage.getItem("token"),
			username: this.$route.params.username,
			comments: {},
			likes: {},
			userLikes: {},
			followTag: true,
			banTag: true,
			newComment: [],
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
				let apiUrl = `/profile/${this.username}/`;
				let response = await this.$axios.get(apiUrl, {
					headers: {
						Authorization:
							"Bearer " + localStorage.getItem("token"),
					},
				});
				this.user = response.data;
				let currentUser = this.user.username;
				let banResponse = await this.$axios.get(
					`/users/${this.userId}/banned/`,
					{
						headers: {
							Authorization:
								"Bearer " + localStorage.getItem("token"),
						},
					}
				);
				let banned = banResponse.data;
				this.banTag = banned
					.map((p) => p.bannedUser)
					.includes(currentUser);

				let follower = localStorage.getItem("username");
				this.followTag = this.user.followers.includes(follower);

				for (const photo of this.user.photosId) {
					const commentsResponse = await this.$axios.get(
						`/photos/${photo}/comments/`,
						{
							headers: {
								Authorization:
									"Bearer " + localStorage.getItem("token"),
							},
						}
					);
					this.comments[photo] = commentsResponse.data;
					const likesResponse = await this.$axios.get(
						`/photos/${photo}/likes/${this.user.id}/`,
						{
							headers: {
								Authorization:
									"Bearer " + localStorage.getItem("token"),
							},
						}
					);
					this.likes[photo] = likesResponse.data.length;
					this.userLikes[photo] = likesResponse.data
						.map((p) => p.userId)
						.includes(parseInt(this.userId));
				}
			} catch (e) {
				if (e.response.status == 401) {
					this.errormsg = "You can not acces" + this.username;
				} else if (e.response.status == 401) {
					// this.$router.push({ name: "Login" });
					this.errormsg = "User not found" + this.username;
				} else {
					this.errormsg = e.toString();
				}
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async likePhoto(photoId, userId) {
			this.loading = true;
			this.errormsg = null;
			try {
				if (this.userLikes[photoId]) {
					await this.$axios.delete(
						`/photos/${photoId}/likes/${userId}/`,
						{
							headers: {
								Authorization:
									"Bearer " + localStorage.getItem("token"),
							},
						}
					);
					this.userLikes[photoId] = false;
				} else {
					await this.$axios.put(
						`/photos/${photoId}/likes/${userId}/`,
						{},
						{
							headers: {
								Authorization:
									"Bearer " + localStorage.getItem("token"),
							},
						}
					);

					this.userLikes[photoId] = true;
				}
				await this.refresh();
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async uncommentPhoto(photoId, commentId) {
			this.loading = true;
			this.errormsg = null;
			try {
				await this.$axios.delete(
					`/photos/${photoId}/comments/${commentId}/`,
					{
						headers: {
							Authorization:
								"Bearer " + localStorage.getItem("token"),
						},
					}
				);
				await this.refresh();
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		commentPhoto: async function (photoId) {
			this.loading = true;
			this.errormsg = null;
			try {
				await this.$axios.post(
					`/photos/${photoId}/comments/`,
					{
						userId: parseInt(this.userId),
						comment: this.newComment[photoId],
					},
					{
						headers: {
							Authorization:
								"Bearer " + localStorage.getItem("token"),
						},
					}
				);

				await this.refresh();
				this.newComment[photoId] = "";
			} catch (e) {
				if (e.response.status == 400) {
					this.errormsg = "Wrong format received";
				} else {
					this.errormsg = e.toString();
				}
			}
			this.loading = false;
		},
		async followersView(id, username) {
			this.$router.push({
				name: "Followers",
				params: { userId: id, username: username },
			});
		},
		async followingsView(id, username) {
			this.$router.push({
				name: "Followings",
				params: { userId: id, username: username },
			});
		},

		async followUser(id, followedUserId) {
			this.loading = true;
			this.errormsg = null;
			try {
				if (this.followTag) {
					await this.$axios.delete(
						`/users/${id}/following/${followedUserId}/`,
						{
							headers: {
								Authorization:
									"Bearer " + localStorage.getItem("token"),
							},
						}
					);
					this.followTag = false;
				} else {
					await this.$axios.put(
						`/users/${id}/following/${followedUserId}/`,
						{},
						{
							headers: {
								Authorization:
									"Bearer " + localStorage.getItem("token"),
							},
						}
					);
					this.followTag = true;
				}
				await this.refresh();
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async banUser(id, bannedUserId) {
			this.loading = true;
			this.errormsg = null;
			try {
				if (this.banTag) {
					await this.$axios.delete(
						`/users/${id}/banned/${bannedUserId}/`,
						{
							headers: {
								Authorization:
									"Bearer " + localStorage.getItem("token"),
							},
						}
					);
					this.banTag = false;
				} else {
					await this.$axios.put(
						`/users/${id}/banned/${bannedUserId}/`,
						{},
						{
							headers: {
								Authorization:
									"Bearer " + localStorage.getItem("token"),
							},
						}
					);
					this.banTag = true;
				}
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		requirePhoto(idPhoto) {
			return `../src/assets/images/${
				this.user.photos[this.user.photosId.indexOf(idPhoto)]
			}`;
		},
		async deletePhoto(photo) {
			this.loading = true;
			this.errormsg = null;
			try {
				await this.$axios.delete(`/photos/${photo}/`, {
					headers: {
						Authorization:
							"Bearer " + localStorage.getItem("token"),
					},
				});
				await this.refresh();
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
			await this.refresh();
		},
	},
	mounted() {
		this.refresh();
	},
};
</script>

<template>
	<div class="background">
		<div class="col-md-3 col-lg-2 empty-col"></div>
		<div class="col-md-9 col-lg-10">
			<div class="d-flex justify-content-center align-items-center h-70">
				<div
					class="card2 bg-white text-dark rounded-5 position-relative"
				>
					<div class="card-body p-5 text-center">
						<h2 class="fw-bold mb-4 text-uppercase">
							@{{ this.user?.username }}
						</h2>
						<div
							v-if="
								parseInt(this.user?.id) ===
								parseInt(this.userId)
							"
						>
							<div
								class="dropdown"
								style="position: absolute; bottom: 0; right: 0"
							>
								<button
									class="btn btn-secondary dropdown-toggle"
									type="button"
									data-bs-toggle="dropdown"
									aria-expanded="false"
								>
									<svg class="feather">
										<use
											href="/feather-sprite-v4.29.0.svg#more-vertical"
										/>
									</svg>
								</button>
								<ul
									class="dropdown-menu"
									aria-labelledby="dropdownMenuButton"
								>
									<li class="nav-item">
										<RouterLink
											:to="{
												name: 'Banned',
												params: { userId: this.userId },
											}"
											class="nav-link"
										>
											Banned Users
										</RouterLink>
									</li>
								</ul>
							</div>
						</div>

						<div class="info">
							<div>
								<a
									href="javascript:"
									class="text-muted mb-1 larger-text"
									style="text-decoration: none"
									@click="
										followersView(
											this.userId,
											this.user?.username
										)
									"
									>Followers</a
								>
								<p class="mb-0 text-center">
									{{ this.user?.followersCount }}
								</p>
							</div>
							<div>
								<a
									href="javascript:"
									class="text-muted mb-1 larger-text"
									style="text-decoration: none"
									@click="
										followingsView(
											this.userId,
											this.user?.username
										)
									"
									>Followings</a
								>
								<p class="mb-0 text-center">
									{{ this.user?.followingCount }}
								</p>
							</div>
							<div>
								<p class="text-muted mb-1 larger-text">
									Photos
								</p>
								<p class="mb-0 text-center">
									{{ this.user?.photosCount }}
								</p>
							</div>
						</div>
						<div
							v-if="
								parseInt(this.user?.id) !==
								parseInt(this.userId)
							"
						>
							<div class="gap-3">
								<button
									class="btn btn-primary btn-block rounded-pill larger-text"
									type="submit"
									@click="banUser(this.userId, this.user?.id)"
									style="background-color: #d10606f5"
								>
									{{ this.banTag ? "Unban" : "Ban" }}
								</button>
								<button
									class="btn btn-primary btn-block rounded-pill larger-text"
									type="submit"
									@click="
										followUser(this.userId, this.user?.id)
									"
									style="background-color: #2e4a78"
								>
									{{ this.followTag ? "Unfollow" : "Follow" }}
								</button>
								<LoadingSpinner v-if="loading" />
							</div>
						</div>
					</div>
				</div>
			</div>
			<div
				v-if="parseInt(this.user?.id) === parseInt(this.userId)"
				class="d-flex justify-content-center align-items-center h-80"
			>
				<router-link
					:to="{
						name: 'UploadPhoto',
					}"
					class="btn btn-primary rounded-pill larger-text"
					style="background-color: #2e4a78"
					>Upload Photo</router-link
				>
				<router-link
					:to="{
						name: 'SetMyUsername',
					}"
					class="btn btn-primary rounded-pill larger-text"
					style="background-color: #2e4a78"
					>Update my Username
				</router-link>
			</div>

			<div
				class="photo-grid d-flex flex-wrap justify-content-center align-items-center h-80"
			>
				<div
					class="col-9 col-sm-6 col-md-4 col-lg-3 mb-5"
					v-for="photo in this.user?.photosId"
					:key="photo.id"
				>
					<div
						class="photo-card bg-white text-dark rounded-5 position-relative"
					>
						<div
							v-if="
								parseInt(this.user?.id) ===
								parseInt(this.userId)
							"
						>
							<div
								class="dropdown position-absolute top-0 end-0 p-2"
							>
								<button
									class="btn btn-secondary dropdown-toggle"
									type="button"
									data-bs-toggle="dropdown"
									aria-expanded="false"
								>
									<svg class="feather">
										<use
											href="/feather-sprite-v4.29.0.svg#more-vertical"
										/>
									</svg>
								</button>
								<ul
									class="dropdown-menu"
									aria-labelledby="dropdownMenuButton"
								>
									<li>
										<button
											class="dropdown-item"
											@click="deletePhoto(photo)"
										>
											Delete Photo
										</button>
									</li>
								</ul>
							</div>
						</div>

						<img :src="requirePhoto(photo)" class="photo-img" />
						<div class="card-body">
							<div
								class="d-flex justify-content-between align-items-center mb-2"
							>
								<div class="likes">
									<button
										class="btn btn-link"
										:class="{
											'text-danger': userLikes[photo],
											'text-dark': !userLikes[photo],
										}"
										@click="likePhoto(photo, this.userId)"
									>
										<svg class="feather">
											<use
												href="/feather-sprite-v4.29.0.svg#heart"
											/>
										</svg>
									</button>
									{{ this.likes[photo] }} likes
								</div>
								<div class="comments" v-if="comments[photo]">
									{{ comments[photo].length }} comments
								</div>
							</div>
							<div class="comment-section">
								<div
									class="comment"
									v-for="comment in comments[photo]"
									:key="comment.id"
								>
									<router-link
										:to="{
											name: 'MyAccount',
											params: {
												username: comment.username,
											},
										}"
										style="
											text-decoration: none;
											color: black;
											font-weight: bold;
										"
									>
										{{ comment.username }}
									</router-link>
									{{ comment.comment }}
									<button
										v-if="
											parseInt(comment.userId) ===
											parseInt(this.userId)
										"
										class="btn btn-link text-danger m1-auto"
										@click="
											uncommentPhoto(photo, comment.id)
										"
									>
										<svg class="feather">
											<use
												href="/feather-sprite-v4.29.0.svg#trash"
											/>
										</svg>
									</button>
								</div>

								<div class="mt-3">
									<input
										type="text"
										v-model="newComment[photo]"
										class="form-control"
										placeholder="Add a comment..."
										@keyup.enter="commentPhoto(photo)"
									/>
								</div>
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
