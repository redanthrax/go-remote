import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import AgentsView from '../views/AgentsView.vue'
import AgentView from '../views/AgentView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
        path: '/',
        name: 'home',
        component: HomeView
    },
    {
        path: '/agents',
        name: 'agents',
        component: AgentsView
    },
    {
        path: '/agent',
        name: 'agent',
        component: AgentView
    }
  ]
})

export default router
