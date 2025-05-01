<script setup lang="ts">
import type { Declaration } from '@/types/Declaration';
import ArrowIcon from '../icons/ArrowIcon.vue';
import SadFaceIcon from '../icons/SadFaceIcon.vue';
const props = defineProps<{
    declaration: Declaration
}>()

</script>

<template>
    <div class="current-declaration">
        <div class="header">
            <div class="title">Текущая декларация</div>
            <div class="text-[24px]">
                <ArrowIcon />
            </div>
        </div>
        <div class="status">
            <div class="progress">{{ declaration.progress }} / 5</div>
            <div class="progress-bar">
                <div class="progress-fill" :style="{ width: (declaration.progress / 5 * 100) + '%' }"></div>
            </div>
            <div class="flex flex-row items-center gap-1">
                <span class="text-white text-[16px]">
                    <SadFaceIcon />
                </span>
                <span class="text-[#FF6644]">грозит -{{ declaration.threat }} ₽</span>
            </div>
        </div>
    </div>
</template>


<style scoped>
.current-declaration {
    background-image: url('../../assets/images/declaration-bg.png');
    background-size: cover;
    background-position: bottom;
    background-repeat: no-repeat;
    position: relative;
    @apply p-4 h-full rounded-[16px] min-w-[189px]
}

.header{
    @apply flex flex-row justify-between
}

.status {
    @apply flex flex-col gap-2 bg-white p-3 rounded-[16px] absolute w-[95%] bottom-[4px] left-1/2 -translate-x-1/2;
}

.progress {
    font-family: "SF Pro Display";
    @apply text-[28px] leading-8
}

.progress-bar {
    @apply h-[24px] rounded-[20px] bg-[#ebebeb] overflow-hidden
}

.progress-fill {
    @apply h-[24px] bg-[#FF6644]
}
</style>

