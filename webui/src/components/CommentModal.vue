<script>
export default {	
	data(){
		return{
            token: localStorage.getItem('token'),
			commentText:"",
		}
	},

	props:['modalID','commentsList','isItMe','photoID'],

	methods: {
		async commentPhoto() {
            try {
                // POST /photos/{pid}/comments/
                let response = await this.$axios.post(`/photos/${this.photoID}/comments/`, this.commentText, {headers: {'Authorization': `${localStorage.getItem('token')}`, 'Content-Type': 'text/plain'}});
                let comment = response.data;
                this.$emit('addComment', comment); // signal to parent
                this.commentText = "";
            } catch(error) {
                const status = error.response.status;
                const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
                alert(this.errormsg);
            }
        },
        async uncommentPhoto(commentID) {
            try {
                // DELETE /photos/{pid}/comments/{cid}
                await this.$axios.delete(`/photos/${this.photoID}/comments/${commentID}`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                this.$emit('removeComment', commentID); // signal to parent
            } catch(error) {
                const status = error.response.status;
                const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
                alert(this.errormsg);
            }
        },

	},
}
</script>

<template>
    <div class="modal fade my-modal-disp-none" :id="modalID" tabindex="-1" aria-hidden="true">
        <div class="modal-dialog modal-dialog-centered modal-dialog modal-dialog-scrollable">
            <div class="modal-content">
                <div class="modal-header">
                    <h1 class="modal-title fs-5" :id="modalID">Comments</h1>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>
                <div v-for="comment in commentsList" :key="comment.commentID">
                    <div class="modal-body">
                        <div class="container-fluid">
                            <hr>
                            <div class="row">
                                <div class="col-10">
                                    <h5>@{{comment.authorUsername}}</h5>
                                </div>
                                <div class="col-2">
                                    <button v-if="token == comment.authorID || isItMe" class="btn my-btn-comm" @click="uncommentPhoto(comment.commentID)">
                                        <i class="fa-regular fa-trash-can my-trash-icon"></i>
                                    </button>
                                </div>
                            </div>
                            <div class="row">
                                <div class="col-12">
                                    {{comment.commentText}}
                                </div>
                            </div>
                            <hr>
                        </div>
                    </div>
                </div>
                <div class="modal-footer d-flex justify-content-center w-100">
                    <div class="row w-100 ">
                        <div class="col-10">
                            <div class="mb-3 me-auto">
                                <textarea class="form-control" id="exampleFormControlTextarea1" 
								placeholder="Add a comment..." rows="1" maxLength="2200" v-model="commentText"></textarea>
                            </div>
                        </div>
                        <div class="col-2 d-flex align-items-center">
                            <button type="button" class="btn btn-primary" 
							@click.prevent="commentPhoto" 
							:disabled="commentText.length < 1 || commentText.length > 2200">
							Send
							</button>
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
.my-btn-comm{
    border: none;
}
.my-btn-comm:hover{
    border: none;
    color: red;
    transform: scale(1.1);
}
</style>
