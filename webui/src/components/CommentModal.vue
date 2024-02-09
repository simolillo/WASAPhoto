<script>
export default {	
	data(){
		return{
			commentText:"",
		}
	},

	props:['modalID','commentsList','photo_owner','photo_id'],

	methods: {
		async commentPhoto() {
            try {
                // POST /photos/{pid}/comments/
                let response = await this.$axios.post(`/photos/${this.photoID}/comments/`, this.commentText, {headers: {'Authorization': `${localStorage.getItem('token')}`, 'Content-Type': 'text/plain'}});
                let comment = response.data;
                this.$emit('addComment', comment); // signal to parent
                this.commentText = ""
            } catch(error) {
                const status = error.response.status;
                const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
                alert(this.errormsg);
            }
        },
        async uncommentPhoto() {
            try {
                // DELETE /photos/{pid}/comments/{cid}
                await this.$axios.delete(`/photos/${this.photoID}/comments/${this.cid}`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                this.$emit('removeComment', cid); // signal to parent
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
        <div class="modal-dialog modal-dialog-centered modal-dialog modal-dialog-scrollable ">
            <div class="modal-content">
                <div class="modal-header">
                    <h1 class="modal-title fs-5" :id="modalID">Comments</h1>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>
                <div v-for="comment in commentsList">
                    <div class="modal-body">
                        <div class="container-fluid">
                            <hr>
                            <div class="row">
                                <div class="col-10">
                                    <h5>@{{author}}</h5>
                                </div>

                                <div class="col-2">
                                    <button v-if="user === author || user === photo_owner" class="btn my-btn-comm" @click="deleteComment">
                                        <i class="fa-regular fa-trash-can my-trash-icon"></i>
                                    </button>
                                </div>

                            </div>

                            <div class="row">
                                <div class="col-12">
                                    {{content}}
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
								placeholder="Add a comment..." rows="1" maxLength="30" v-model="commentText"></textarea>
                            </div>
                        </div>

                        <div class="col-2 d-flex align-items-center">
                            <button type="button" class="btn btn-primary" 
							@click.prevent="addComment" 
							:disabled="commentText.length < 1 || commentText.length > 30">
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
    color: var(--color-red-danger);
    transform: scale(1.1);
}
</style>
