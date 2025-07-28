<script lang="ts" setup>
import { Declaration } from '@/types/Declaration';
import ArrowIcon from '@/components/icons/ArrowIcon.vue';
import { computed, onBeforeMount } from 'vue';
import HappyFaceIcon from '@/components/icons/HappyFaceIcon.vue'
import SadFaceIcon from '@/components/icons/SadFaceIcon.vue';
import TimeIcon from '@/components/icons/TimeIcon.vue';
import { DeclarationDocument } from '../entities/DeclarationDocument';
import { DeclarationStatus } from '../lib/constants';
import { decl } from 'postcss';
import { useRouter } from 'vue-router';
import { Member } from '@/types/Member';
const props = defineProps<{
    declarations: DeclarationDocument[],
    member: Member
}>()
const router = useRouter()

function formatDate(dateStr:string) {
    if (dateStr == "") return "Не установлена"
    const date = new Date(dateStr)
    const day = date.getUTCDate().toString().padStart(2, "0");
    const month = (date.getMonth() + 1).toString().padStart(2, "0")
    const year = date.getFullYear();

    return `${day}.${month}.${year}`
}

function declarationStatus(declaration: DeclarationDocument):string {
    const declaration_date = new Date(declaration.end_date? declaration.end_date : "")
    const now = new Date();
    if (now < declaration_date) {
        return "progress"
    }

    return "ok"
}

function getBgColor(declaration:DeclarationDocument):string{
    switch(declaration.status) {
        case DeclarationStatus.inProgress || DeclarationStatus.notStared:
            return "black"
        case DeclarationStatus.failed:
            return "#efd0ca"
        case DeclarationStatus.done:
            return "#c5e2cd"
    }
    return ''
}

function getIcon(declaration:DeclarationDocument) {
    switch(declaration.status) {
        case DeclarationStatus.inProgress || DeclarationStatus.notStared:
            return TimeIcon
        case DeclarationStatus.failed:
            return SadFaceIcon
        case DeclarationStatus.done:
            return HappyFaceIcon
    }
    return SadFaceIcon
}

function getStatusText(declaration:DeclarationDocument) {
    const end_date = new Date(declaration.end_date? declaration.end_date : "")
    const now = new Date()
    if (now < end_date) {
        return "Текущая"
    }
    return "Прошлая"
}

function openDeclaration(declaration: DeclarationDocument) {
    router.push(`/declaration/${props.member.username}/${declaration.id}`)
}

onBeforeMount(() => {
    console.log("DECLARATIONS: ", props.declarations)
})

</script>

<template>
    <div class="declarations">
        <div class="title">Декларации</div>
      <div
        v-if="props.declarations && props.declarations.length > 0"
        class="single-declaration"
        v-for="declaration in declarations"
        :key="declaration.id"
        :style="{ backgroundColor: getBgColor(declaration) }"
        @click="() => openDeclaration(declaration)"
      >
        <div class="info" :class="declarationStatus(declaration)" :style="{color: declarationStatus(declaration) == 'progress'? 'white': 'black'}">
          <component class="icon" :is="getIcon(declaration)" :style="{color: getBgColor(declaration)}" />
          <span>{{ getStatusText(declaration) }}</span>
          <span>{{ `До ${formatDate(declaration.end_date? declaration.end_date : "")}` }}</span>
        </div>
        <div class="link-sign" :style="{color: declarationStatus(declaration) == 'progress'? 'white': 'black'}">
          <ArrowIcon />
        </div>
      </div>
    </div>
  </template>

<style scoped>

.declarations {
    @apply flex flex-col p-[6px] bg-white rounded-[20px] gap-1
}

.title {
    font-family: "SF Pro Display";
    @apply pt-3 pr-4 pb-3 pl-4 font-[500] text-[28px] leading-8 tracking-[-0.4]
}

.link-sign{
    color:white;
}

.single-declaration {
    @apply flex flex-row justify-between rounded-[16px] p-4

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

.empty-message {
    font-family: "SF Pro Text";
    @apply text-left text-base font-medium w-full
}
</style>