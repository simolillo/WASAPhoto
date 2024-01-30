import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: HomeView},
		{path: '/link1', component: HomeView},
		{path: '/login', component: LoginView},
		{path: '/some/:id/link', component: HomeView},
	]
})

export default router
