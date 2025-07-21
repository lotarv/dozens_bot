<script lang="ts" setup>
import ArrowLeft from '@/components/icons/ArrowLeft.vue';
import { useBankStore } from '../model/bankStore';
import CrossIcon from '../ui/icons/CrossIcon.vue';
import { onMounted } from 'vue';
import { formatDateTime } from '@/services/helpers';
const baseURL = import.meta.env.VITE_BASE_URL
const bankStore = useBankStore()

onMounted(() => {
    console.log(bankStore.bank)
})
</script>

<template>
    <main class="flex flex-col gap-[10px]">
        <div class="header">
            <RouterLink :to="'/'">
                <div class="p-2 bg-white w-12 h-12 rounded-full text-[22px] flex items-center justify-center"><ArrowLeft/></div>
            </RouterLink>
            <span>Копилка</span>
        </div>
        <div class="w-full bg-white flex flex-row rounded-2xl p-4 justify-between items-center">
            <div class="flex flex-col gap-2">
                <span class="font-medium text-[28px] leading-7 tracking-[-0.4px]">{{ bankStore.bank?.balance.toLocaleString("ru-RU") }} ₽</span>
                <span class="font-semibold text-base tracking-[-0.4px]">Сейчас в копилке</span>
            </div>
            <div class="rounded-full bg-[#EBEBEB] flex items-center justify-center h-12 w-12">
                <CrossIcon class="text-3xl"></CrossIcon>
            </div>
        </div>
        <div class="history">
            <p class="font-medium text-[28px] leading-8 tracking-[-0.4px] py-3 px-4">История</p>
            <div class="transactions-container">
                <div v-for="transaction of bankStore.bank?.transactions" class="transaction">
                    <div class="member-info flex gap-2">
                        <img :src="`${transaction.member.avatar_url}`" alt="" class="rounded-full w-12 h-12">
                        <div class="flex flex-col gap-1">
                            <span class="font-medium text-base tracking-[-0.4px]">{{ transaction.member.full_name }}</span>
                            <span class="text-[#767676] text-sm tracking-[-0.25px]">{{ transaction.reason }}</span>
                        </div>
                    </div>
                    <div class="transaction-info flex flex-col gap-1">
                        <span class="text-[#2EC055] font-semibold text-base" v-if="transaction.amount > 0">+ {{ transaction.amount.toLocaleString("ru-RU")  }} ₽</span>
                        <span class="text-[#FF6644] font-semibold text-base" v-else>- {{ Math.abs(transaction.amount).toLocaleString("ru-RU") }} ₽</span>
                        <span class="text-[#767676] text-sm tracking-[-0.25px]">{{ formatDateTime(transaction.created_at) }}</span>
                    </div>
                </div>
            </div>
        </div>
    </main>
</template>
<style scoped>
.header {
    @apply text-[36px] font-medium tracking-[-1px] flex items-center p-3 gap-4
}
.transaction {
    @apply flex justify-between p-4 border-b border-gray-300
}

.transactions-container>*:last-child {
  @apply border-b-0;
}


</style>