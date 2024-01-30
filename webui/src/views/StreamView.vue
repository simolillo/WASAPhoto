<script>
import "../styles.css";

export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			userId: localStorage.getItem("token"),
			photos: [],
			likes: {},
			comments: {},
			userLikes: {},
			newComment: null,
		};
	},
	methods: {
		load() {
			return load;
		},
		requirePhoto(path) {
			return `../src/assets/images/${path}`;
		},
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let apiUrl = `/users/${this.userId}/stream/`;
				let response = await this.$axios.get(apiUrl, {
					headers: {
						Authorization:
							"Bearer " + localStorage.getItem("token"),
					},
				});
				this.photos = response.data;
				for (const photo of this.photos) {
					const commentsResponse = await this.$axios.get(
						`/photos/${photo.id}/comments/`,
						{
							headers: {
								Authorization:
									"Bearer " + localStorage.getItem("token"),
							},
						}
					);
					this.comments[photo.id] = commentsResponse.data;
					const likesResponse = await this.$axios.get(
						`/photos/${photo.id}/likes/${photo.userId}/`,
						{
							headers: {
								Authorization:
									"Bearer " + localStorage.getItem("token"),
							},
						}
					);
					this.likes[photo.id] = likesResponse.data.length;
					this.userLikes[photo.id] = likesResponse.data
						.map((p) => p.userId)
						.includes(parseInt(this.userId));
				}
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async likePhoto(photoId, id) {
			this.loading = true;
			this.errormsg = null;
			try {
				if (this.userLikes[photoId]) {
					await this.$axios.delete(
						`/photos/${photoId}/likes/${id}/`,
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
						`/photos/${photoId}/likes/${id}/`,
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
		async deleteComment(photoId, commentId) {
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
		addComment: async function (photoId) {
			this.loading = true;
			this.errormsg = null;
			try {
				await this.$axios.post(
					`/photos/${photoId}/comments/`,
					{
						userId: parseInt(this.userId),
						comment: this.newComment,
					},
					{
						headers: {
							Authorization:
								"Bearer " + localStorage.getItem("token"),
						},
					}
				);
				await this.refresh();
				this.newComment = "";
			} catch (e) {
				if (e.response.status == 400) {
					this.errormsg = "Wrong format received " + this.newComment;
				} else {
					this.errormsg = e.toString();
				}
			}
			this.loading = false;
		},
	},

	mounted() {
		this.refresh();
	},
};
</script>

<template>
	<div class="background">
		<div class="container py-5">
			<div class="d-flex justify-content-center align-items-center h-80">
				<router-link
					:to="{
						name: 'SearchUser',
					}"
					class="btn btn-primary rounded-pill larger-text"
					style="background-color: #2e4a78"
					>Search User</router-link
				>
			</div>
			<div class="row justify-content-center align-items-start">
				<div class="col-md-5">
					<div
						class="photo-card bg-white text-dark rounded-5 mb-4"
						v-for="photo in photos"
						:key="photo.id"
					>
						<div class="position-relative">
							<img
								:src="requirePhoto(photo.path)"
								class="photo-img"
							/>
							<div
								class="photo-overlay position-absolute top-0 start-0 p-2"
							>
								<router-link
									:to="{
										name: 'MyAccount',
										params: { username: photo.username },
									}"
									style="
										text-decoration: none;
										color: white;
										font-weight: bold;
									"
								>
									{{ photo.username }}
								</router-link>
								<span style="margin-left: 5px">
									{{
										new Date(photo.DateTime).toLocaleString(
											"en-GB",
											{
												dateStyle: "short",
												timeStyle: "short",
											}
										)
									}}</span
								>
							</div>
						</div>
						<div class="card-body">
							<div
								class="d-flex justify-content-between align-items-center mb-2"
							>
								<div class="likes">
									<button
										class="btn btn-link"
										:class="{
											'text-danger': userLikes[photo.id],
											'text-dark': !userLikes[photo.id],
										}"
										@click="
											likePhoto(photo.id, this.userId)
										"
									>
										<svg class="feather">
											<use
												href="/feather-sprite-v4.29.0.svg#heart"
											/>
										</svg>
									</button>
									{{ this.likes[photo.id] }} likes
								</div>
								<div class="comments" v-if="comments[photo.id]">
									{{ comments[photo.id].length }} comments
								</div>
							</div>
							<div class="comment-section">
								<div
									class="comment"
									v-for="comment in comments[photo.id]"
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
											deleteComment(photo.id, comment.id)
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
										v-model="newComment"
										class="form-control"
										placeholder="Add a comment..."
										@keyup.enter="addComment(photo.id)"
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
