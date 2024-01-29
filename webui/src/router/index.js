import {createRouter, createWebHashHistory} from 'vue-router'
import LogInView from '../views/LogInView.vue'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		// {path: '/', redirect: '/login'},
		// {path: '/login', component: LogInView},
		{path: '/', component: HomeView},
		{path: '/link1', component: HomeView},
		{path: '/link2', component: HomeView},
		{path: '/some/:id/link', component: HomeView},
	]
})

export default router
