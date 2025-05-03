import {createRouter, createWebHistory} from 'vue-router'
const routes = [
    {
        path: '/',
        name:'MainView',
        component: () => import('../views/MainView.vue')
    },
    {
        path: '/members',
        name:'members',
        component: () => import('../views/MembersView.vue')
    },
    {
        path: '/declarations/:username',
        name: "declaration",
        component: () => import('../views/DeclarationsView.vue'),
        props:true
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes,
}) 

export default router