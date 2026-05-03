/// <reference types="vite/client" />
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/:pathMatch(.*)*',
            name: 'dynamic-page',
            component: () => import('../pages/DynamicPage.vue')
        }
    ]
})

export default router
