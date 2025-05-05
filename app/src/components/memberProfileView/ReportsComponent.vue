<script lang="ts" setup>
import { Report } from '@/types/Report';
import ArrowIcon from '../icons/ArrowIcon.vue';
import { computed, ref } from 'vue'
const props = defineProps<{
    reports: Report[],
}>()

const visibleCount = ref<number>(6)

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
function formatDate(dateStr:string) {
    const date = new Date(dateStr)
    const day = date.getUTCDate();
    const month = date.getMonth() + 1
    const year = date.getFullYear();

    return `${day}.${month}.${year}`
}
</script>
<template>
    <div class="reports">
        <div class="title">Отчеты</div>
        <div class="flex-container">
            <div class="single-report" v-for="report in visibleReports">
                <div class="report-info">
                    <span class="date">{{ formatDate(report.date) }}</span>
                    <span class="text-[24px]"><ArrowIcon/></span>
                </div>
            </div>
        </div>
        <div v-if="hasMoreReports" class="loadMoreButton" @click="loadMore()">Смотреть еще</div>
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

.flex-container{
    @apply flex flex-wrap justify-between gap-1
}

.single-report {
    font-family: "SF Pro Text";
    @apply p-4 rounded-[16px] bg-white min-w-[185px] flex-1 font-[500] leading-6 tracking-[-0.4px]
}

.report-info{
    @apply flex justify-between items-center 
}

.loadMoreButton {
    font-family: "SF Pro Text";
    @apply font-bold text-base leading-6 tracking-[-0.4px] text-center p-4 cursor-pointer
}
</style>