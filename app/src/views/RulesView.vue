<script lang="ts" setup>
import ArrowLeft from '@/components/icons/ArrowLeft.vue';
import MarkDownComponent from '@/components/MarkDownComponent.vue';
import { getTelegramInitData } from '@/services/auth';
import { decryptGoAES, isEncrypted } from '@/services/crypto';
import axios from 'axios';
import { onMounted, ref } from 'vue';
import { useDecryptionStore } from '@/stores/decryption';

const cryptoStore = useDecryptionStore()
const isLoading = ref(false)
const error = ref<string | null>(null)
const rulesStr = ref<string>('')
onMounted(async () => {
    try {
        isLoading.value = true
        error.value = null
        const response = await axios.get(`${import.meta.env.VITE_API_URL}/rules`, {
            headers: {
                'X-Telegram-Init-Data': getTelegramInitData(),
                "Cache-Control":"no-cache"
            }
        })
        rulesStr.value = response.data.text

    } catch (err) {
        error.value = 'Failed to load rules. Please try again later.';
        console.error('Failed to load rules:', err);
    } finally {
        isLoading.value = false
    }
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
            <span>Правила</span>
        </div>
        <div v-if="isLoading" class="loading"></div>
        <div v-else class="px-3">
            <MarkDownComponent :text="rulesStr" />
        </div>
    </section>
</template>

<style scoped>
.header {
    @apply text-[36px] font-medium tracking-[-1px] flex items-center p-3 gap-4
}
</style>