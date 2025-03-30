import { createRouter, createWebHistory } from 'vue-router'
import Playground from '../components/Playground.vue'

const routes = [
  {
    path: '/',
    name: 'Playground',
    component: Playground
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router 