<script lang="ts" setup>

import { ref } from 'vue';
import { onBeforeMount } from 'vue';
import { DeclarationDocument } from '../entities/DeclarationDocument';
import { DozensTransport } from '@/repository/http';
import { useRoute } from 'vue-router';
import { useRouter } from 'vue-router';
import MarkDownComponent from '@/components/MarkDownComponent.vue';
import ArrowLeft from '@/components/icons/ArrowLeft.vue';
const route = useRoute()
const router = useRouter()
const declaration = ref<DeclarationDocument | null>(null)

function goBack() {
    router.back()
}

onBeforeMount(async() => {
    const declarationID = route.params.id?.toString()
    if (!declarationID) return

    declaration.value = await DozensTransport.GetDeclarationByID(declarationID)
    console.log(declarationID)
    console.log(declaration)
})
</script>
<template>
      <section class="p-1 flex flex-col font-medium" v-if="declaration !== null">
        <div class="header">
            <div @click="goBack"class="p-2 bg-white w-12 h-12 rounded-full text-[22px] flex items-center justify-center">
                <ArrowLeft />
            </div>
            <p class="flex-1 flex justify-between items-center"><span>Декларация</span> </p>
        </div>
        <div class="px-2">
            <MarkDownComponent :text="declaration.text"></MarkDownComponent>
        </div>
    </section>
</template>

<style scoped>
.header {
    @apply text-[36px] font-medium tracking-[-1px] flex items-center p-3 gap-4
}
</style>