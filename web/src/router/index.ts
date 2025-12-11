import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import ErrorView from '@/views/ErrorView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/error/:status',
      name: 'error',
      component: ErrorView,
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/error/404',
    },
  ],
})

export default router
