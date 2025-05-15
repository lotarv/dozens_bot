<script lang="ts" setup>
import ArrowLeft from '@/components/icons/ArrowLeft.vue';
import { onBeforeMount } from 'vue';
import { useReportsStore } from '@/stores/reportsStore';
import { useUserStore } from '@/stores/user';
import { useRouter } from 'vue-router';
const reportsStore = useReportsStore()
const userStore = useUserStore()
const router = useRouter()
function formatDate(isoDate: string): string {
    const date = new Date(isoDate)
    const day = String(date.getDate()).padStart(2, "0")
    const month = String(date.getMonth() + 1).padStart(2, "0") // месяцы от 0
    const year = date.getFullYear()

    return `${day}.${month}.${year}`
}

function openReport(username: string, reportID: number) {
    router.push({name: "report", params: {username: username, id:reportID }})
}
onBeforeMount(async () => {
    reportsStore.fetchUserReports(userStore.currentUser.username)
})

</script>
<template>
    <section class="p-1 flex flex-col font-medium">
        <div class="header">
            <RouterLink :to="'/'">
                <div class="p-2 bg-white w-12 h-12 rounded-full text-[22px] flex items-center justify-center">
                    <ArrowLeft />
                </div>
            </RouterLink>
            <span>Отчеты</span>
        </div>

        <div class="reports-container">
            <div @click="openReport(userStore.currentUser.username,index)" v-for="(report, index) in reportsStore.reports[userStore.currentUser.username]?.reports" class="report">
                <img :src="userStore.currentUser.avatar_url" alt="" class="avatar">
                <div class="date">
                    <span class="date-title">Дата</span>
                    <span class="date-date">{{ formatDate(report.creation_date) }}</span>
                </div>
            </div>
        </div>
    </section>
</template>

<style scoped>
.header {
    @apply text-[36px] font-medium tracking-[-1px] flex items-center p-3 gap-4
}

.reports-container{
    @apply flex flex-col gap-1
}

.report {
    @apply bg-white flex justify-start items-center gap-3 p-4 rounded-[16px]
}

.avatar{
    @apply w-12 h-12 rounded-full
}

.date{
    font-family: "SF Pro Text";
    @apply flex flex-col text-base leading-6 tracking-[-0.4px]
}

.date-title {
    @apply text-gray-500
}
</style>
