<script lang="ts" setup>
import { Declaration } from '@/types/Declaration';
import ArrowIcon from '../icons/ArrowIcon.vue';
import { computed } from 'vue';
import HappyFaceIcon from '../icons/HappyFaceIcon.vue';
import SadFaceIcon from '../icons/SadFaceIcon.vue';
import TimeIcon from '../icons/TimeIcon.vue';
const props = defineProps<{
    declarations: Declaration[],
}>()

function formatDate(dateStr:string) {
    const date = new Date(dateStr)
    const day = date.getUTCDate();
    const month = date.getMonth() + 1
    const year = date.getFullYear();

    return `${day}.${month}.${year}`
}

function declarationStatus(declaration: Declaration):string {
    const declaration_date = new Date(declaration.date)
    const now = new Date();
    if (now < declaration_date) {
        return "progress"
    }

    if (declaration.progress != 5) {
        return "fail"
    }

    return "ok"
}

function getBgColor(status: string):string{
    switch(status) {
        case "progress":
            return "black"
        case "fail":
            return "#efd0ca"
        case "ok":
            return "#c5e2cd"
    }
    return ''
}

function getIcon(status: string) {
    switch(status) {
        case "progress":
            return TimeIcon
        case "fail":
            return SadFaceIcon
        case "ok":
            return HappyFaceIcon
    }
    return SadFaceIcon
}

function getStatusText(status:string) {
    if (status == "progress") {
        return "Текущая."
    } 
    return "Прошлая."
}


</script>

<template>
    <div class="declarations">
        <div class="title">Декларации</div>
      <div
        class="single-declaration"
        v-for="declaration in declarations"
        :key="declaration.date"
        :style="{ backgroundColor: getBgColor(declarationStatus(declaration)) }"
      >
        <div class="info" :class="declarationStatus(declaration)" :style="{color: declarationStatus(declaration) == 'progress'? 'white': 'black'}">
          <component class="icon" :is="getIcon(declarationStatus(declaration))" :style="{color: getBgColor(declarationStatus(declaration))}" />
          <span>{{ getStatusText(declarationStatus(declaration)) }}</span>
          <span>{{ `До ${formatDate(declaration.date)}` }}</span>
        </div>
        <div class="link-sign" :style="{color: declarationStatus(declaration) == 'progress'? 'white': 'black'}">
          <ArrowIcon />
        </div>
      </div>
    </div>
  </template>

<style scoped>

.declarations {
    @apply flex flex-col p-2 bg-white rounded-[24px] gap-1
}

.title {
    font-family: "SF Pro Display";
    @apply pt-3 pr-4 pb-3 pl-4 font-[500] text-[28px] leading-8 tracking-[-0.4]
}

.link-sign{
    color:white;
}

.single-declaration {
    @apply flex flex-row justify-between rounded-[16px] p-4 opacity-[0.55]

}

.info{
    font-family: "SF Pro Text";
    @apply flex gap-1 items-center font-[500] text-base tracking-[-0.4] leading-6
}

.icon{
    @apply text-[20px] mr-1
}

.link-sign {
    @apply text-[24px]
}
</style>