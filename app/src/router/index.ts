import {createRouter, createWebHistory} from 'vue-router'
import MainView from '../views/MainView.vue'
import MembersView from '../views/MembersView.vue'
import MemberProfileView from '@/views/MemberProfileView.vue'
import RulesView from '@/views/RulesView.vue'
import ReportsView from '@/views/ReportsView.vue'
import ReportView from '@/views/ReportView.vue'
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
    },
    {
        path:"/rules",
        name:"rules",
        component: RulesView
    },
    {
        path:"/reports",
        name:"reports",
        component: ReportsView
    },
    {
        path:"/report/:username/:id",
        name:"report",
        component: ReportView
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes,
}) 

export default router