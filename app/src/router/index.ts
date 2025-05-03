import {createRouter, createWebHistory} from 'vue-router'
import MainView from '../views/MainView.vue'
import MembersView from '../views/MembersView.vue'
import DeclarationsView from '../views/DeclarationsView.vue'
const routes = [
    {
        path: '/',
        name:'MainView',
        component: MainView
    },
    {
        path: '/members',
        name:'members',
        component: MembersView
    },
    {
        path: '/declarations/:username',
        name: "declaration",
        component: DeclarationsView,
        props:true
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes,
}) 

export default router