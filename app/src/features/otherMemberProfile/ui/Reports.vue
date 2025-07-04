<script lang="ts" setup>
import { Report } from '@/types/Report';
import ArrowIcon from '@/components/icons/ArrowIcon.vue';
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router';
const props = defineProps<{
    reports: Report[],
    username:string,
}>()

const visibleCount = ref<number>(6)
const router = useRouter()
const visibleReports = computed(() => {
    return props.reports.slice(0, visibleCount.value)
})

const hasMoreReports = computed(() => {
    return visibleCount.value < props.reports.length
})

function loadMore() {
    if (hasMoreReports.value) {
        visibleCount.value += 6
    }
}
function formatDate(dateStr: string) {
    const date = new Date(dateStr)
    const day = date.getUTCDate().toString().padStart(2, '0');
    const month = (date.getMonth() + 1).toString().padStart(2, '0')
    const year = date.getFullYear();

    return `${day}.${month}.${year}`
}

function openReport(username: string, reportID: number) {
    router.push({name: "report", params: {username: username, id:reportID }})
}
</script>
<template>
    <div class="reports">
        <div class="title">Отчеты</div>

        <div v-if="props.reports.length === 0" class="empty-message">
            У этого участника пока нет отчётов
        </div>

        <div v-else class="flex-container">
            <div @click="openReport(username, index)" class="single-report" v-for="(report, index) in visibleReports" :key="index">
                <div class="report-info">
                    <span class="date">{{ formatDate(report.creation_date) }}</span>
                    <span class="text-[24px]">
                        <ArrowIcon />
                    </span>
                </div>
            </div>
        </div>

        <div v-if="hasMoreReports" class="loadMoreButton" @click="loadMore()">
            Смотреть еще
        </div>
    </div>
</template>

<style scoped>
.reports {
    @apply flex flex-col pt-1 pl-1 pr-1 pb-4 rounded-[20px] gap-1 items-center
}

.title {
    font-family: "SF Pro Display";
    @apply pt-3 pr-4 pb-3 pl-4 font-[500] text-[28px] leading-8 tracking-[-0.4] w-full
}

.empty-message {
    font-family: "SF Pro Text";
    @apply text-left text-base font-medium w-full
}

.flex-container {
    @apply grid grid-cols-2 gap-1 w-full
}

.single-report {
    font-family: "SF Pro Text";
    @apply p-4 rounded-[16px] bg-white flex-1 font-[500] leading-6 tracking-[-0.4px]
}

.report-info {
    @apply flex justify-between items-center
}

.loadMoreButton {
    font-family: "SF Pro Text";
    @apply font-bold text-base leading-6 tracking-[-0.4px] text-center p-4 cursor-pointer
}
</style>