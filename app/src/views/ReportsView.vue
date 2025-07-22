<script lang="ts" setup>
import ArrowLeft from '@/components/icons/ArrowLeft.vue';
import { computed, onBeforeMount, ref } from 'vue';
import { useReportsStore } from '@/stores/reportsStore';
import { useUserStore } from '@/stores/user';
import { useRouter } from 'vue-router';
import { useMembersStore } from '@/stores/membersStore';
const reportsStore = useReportsStore()
const membersStore = useMembersStore()
const userStore = useUserStore()
const router = useRouter()

const usernames = ["annbergen",
    "Dkokorev",
    "MikhailStep",
    "newhndrx",
    "av_yanovskaya",
    "Maxim_Bortsov_Realty",
    "julie_dontsova",
    "polina_happydays",
    "yurii_happy_math",
    "Oldinfree",
    "IL47RUS",
    "incetro"
]

const allReports = ref<any[]>([])
const isReady = ref(false)
onBeforeMount(async () => {
  for (const username of usernames) {
    await reportsStore.fetchUserReports(username)
    const userReports = reportsStore.reports[username]?.reports || []
    const member = membersStore.getMemberByUsername(username)

    userReports.forEach((report, index) => {
      allReports.value.push({
        ...report,
        username,
        fio: member?.fio ?? username,
        avatar_url: member?.avatar_url ?? '',
        dateKey: formatDate(report.creation_date),
        reportIndex: index, // вот здесь сохраняем индекс для перехода
        memberNiche: member?.niche
      })
    })
  }
  isReady.value = true
})

// Группировка по дате
const groupedReports = computed(() => {
  const groups: Record<string, any[]> = {}

  for (const report of allReports.value) {
    const dateKey = new Date(report.creation_date).toISOString().split("T")[0] // => "2025-09-07"
    if (!groups[dateKey]) groups[dateKey] = []
    groups[dateKey].push(report)
  }

  // Сортируем отчёты внутри каждой группы
  for (const key in groups) {
    groups[key].sort((a, b) => new Date(b.creation_date).getTime() - new Date(a.creation_date).getTime())
  }

  // Сортируем группы по дате
  return Object.entries(groups).sort(
    (a, b) => new Date(b[0]).getTime() - new Date(a[0]).getTime()
  )
})

function formatDate(isoDate: string): string {
  const date = new Date(isoDate)
  const day = String(date.getDate()).padStart(2, "0")
  const month = String(date.getMonth() + 1).padStart(2, "0")
  const year = date.getFullYear()
  return `${day}.${month}.${year}`
}

function formatDateTime(isoDate: string): string {
  const date = new Date(isoDate)
  const months = [
    "января", "февраля", "марта", "апреля", "мая", "июня",
    "июля", "августа", "сентября", "октября", "ноября", "декабря"
  ]
  return `${date.getUTCDate()} ${months[date.getUTCMonth()]}, ${date.getUTCFullYear()}`
}

function openReport(username: string, reportID: number) {
  router.push({ name: "report", params: { username, id: reportID } })
}

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

        <div class="reports-container" v-if="isReady">
        <div v-for="[date, reports] in groupedReports" :key="date" class="date-group flex flex-col gap-2">
            <div class="text-lg px-2 pt-4 pb-1">{{ formatDateTime(date) }}</div>
            <div
            v-for="report in reports"
            :key="report.username + '_' + report.reportIndex"
            class="report"
            @click="openReport(report.username, report.reportIndex)"
            >
            <img :src="report.avatar_url" alt="" class="avatar">
            <div class="date">
                <span class="date-title">{{ report.fio }}</span>
                <span class="date-date">{{ report.memberNiche }}</span>
            </div>
            </div>
        </div>
        </div>
        <div v-else class="flex flex-col items-center justify-center">
            <img src="../assets/images/logo.png" alt="" class="h-12 w-12">
            <span>Загрузка отчетов...</span>
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
