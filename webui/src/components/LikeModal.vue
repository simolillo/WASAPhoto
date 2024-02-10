<script>
export default {	

	props:['modalID','likesList'],

    methods: {
        visitAuthorProfile() {
            // /profiles/:username
			this.$router.push(`/profiles/${this.authorUsername}`);
		},
    }
}
</script>

<template>
	<div class="modal fade modal-dialog modal-dialog-centered modal-dialog-scrollable my-modal-disp-none" :id="modal_id" tabindex="-1" aria-hidden="true">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <h1 class="modal-title fs-5" :id="modalID">Likes</h1>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>
                <div class="modal-body">
                    <div v-for="user in likesList">
                        <div class="row mb-2 mt-2">
                            <div class="col d-flex justify-content-center">
                                <div class="user-mini-card card bg-transparent border-start">
                                    <div class="card-body">
                                        <h5 @click="visitAuthorProfile" class="user-mini-card-title d-flex justify-content-center ">@{{ user.userID }}</h5>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
	</div>
</template>

<style>
.my-modal-disp-none{
	display: none;
}
</style>
