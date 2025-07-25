<script lang="ts" setup>
import { useReportsStore } from '@/stores/reportsStore';
import { useDecryptionStore } from '@/stores/decryption';
import { useRoute, useRouter } from 'vue-router';
import MarkDownComponent from '@/components/MarkDownComponent.vue';
import { computed, onBeforeMount, ref } from 'vue';
import ArrowLeft from '@/components/icons/ArrowLeft.vue';
import { decryptGoAES, isEncrypted } from '@/services/crypto';
const router = useRouter()
const route = useRoute()
const reportsStore = useReportsStore()
const cryptoStore = useDecryptionStore()

const username = route.params.username as string
const reportId = Number(route.params.id)

function formatDate(isoDate: string): string {
    const date = new Date(isoDate)
    const day = String(date.getDate()).padStart(2, "0")
    const month = String(date.getMonth() + 1).padStart(2, "0") // месяцы от 0
    const year = date.getFullYear()

    return `${day}.${month}.${year}`
}

const report = computed(() => {
    return reportsStore.reports[username]?.reports[reportId]
})

const reportText = ref<string>("")

function goBack() {
    router.back()
}

onBeforeMount(async() => {
    if (isEncrypted(report.value.text)) {
        reportText.value = await decryptGoAES(report.value.text, cryptoStore.key)
    } else {
        reportText.value = report.value.text // оставить как есть
    }
})

</script>
<template>
    <section class="p-1 flex flex-col font-medium">
        <div class="header">
            <div @click="goBack"class="p-2 bg-white w-12 h-12 rounded-full text-[22px] flex items-center justify-center">
                <ArrowLeft />
            </div>
            <p class="flex-1 flex justify-between items-center"><span>Отчет</span></p>
        </div>
        <div class="px-2">
            <span class="text-base text-gray-700">{{ formatDate(report.creation_date) }}</span>
            <MarkDownComponent :text="reportText"></MarkDownComponent>
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
    @apply bg-white flex justify-start items-center gap-3 p-4 rounded-[16px] cursor-pointer
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
