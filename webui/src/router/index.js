import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import Login from '../views/LoginView.vue'
import StreamView from '../views/StreamView.vue'
import MyAccount from '../views/MyAccount.vue'
import Followings from '../views/GetFollowings.vue'
import Banned from '../views/GetBans.vue'
import Followers from '../views/GetFollowers.vue'
import UploadPhoto from '../views/UploadPhotoView.vue'
import SearchUser from '../views/SearchUser.vue'
import SetMyUsername from '../views/SetMyUsername.vue'


const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{ path: '/', component: Login, name: 'Login' },
		{ path: '/users/stream', component: StreamView, name: 'Stream' },
		{ path: '/users/photos', component: UploadPhoto, name: 'UploadPhoto' },
		{ path: '/profile/:username', component: MyAccount, name: 'MyAccount'},
		{ path: '/users/:userId/followings/:username', component: Followings, name: 'Followings'},
		{ path: '/users/:userId/followers/:username', component: Followers, name: 'Followers'},
		{ path: '/users/:userId/banned', component: Banned, name: 'Banned'},
		{ path: '/search', component: SearchUser, name: 'SearchUser'},
		{ path: '/users/', component: SetMyUsername, name: 'SetMyUsername'}
		
	]
})

export default router
