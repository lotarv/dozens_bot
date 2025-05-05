import {createRouter, createWebHistory} from 'vue-router'
import MainView from '../views/MainView.vue'
import MembersView from '../views/MembersView.vue'
import MemberProfileView from '@/views/MemberProfileView.vue'
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
        path: '/member/:username',
        name: "declaration",
        component: MemberProfileView,
        props:true
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes,
}) 

export default router