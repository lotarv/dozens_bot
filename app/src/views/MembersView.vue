<script lang="ts" setup>
import axios from 'axios';
import { ref, onBeforeMount, Ref } from 'vue';
import { getTelegramInitData } from '../services/auth';
import ArrowIcon from '@/components/icons/ArrowIcon.vue';
import {RouterLink} from 'vue-router'
import { useMembersStore } from '@/stores/membersStore';
import ArrowLeft from '@/components/icons/ArrowLeft.vue';
const membersStore = useMembersStore()
// const members: Ref<Member[]> = ref([]);
const isLoading = membersStore.isLoading
const error = membersStore.error
console.log("LALALALALAL")
onBeforeMount(async () => {
    membersStore.fetchMembers()
})
</script>

<template>
    <section class="p-1 flex flex-col font-medium">
        <div class="header">
            <RouterLink :to="'/'">
                <div class="p-2 bg-white w-12 h-12 rounded-full text-[22px] flex items-center justify-center"><ArrowLeft/></div>
            </RouterLink>
            <span>Участники</span>
        </div>
        <div v-if="isLoading" class="text-center text-xl">Загрузка...</div>
        <div v-else-if="error" class="error">{{ error }}</div>
        <div v-else-if="membersStore.members.length > 0">
            <div class="members-container">
                <RouterLink
                :to="`/member/${member.username}`"
                class="member-card" v-for="member in membersStore.members" :key="member.username">
                    <div class="member-header">
                        <div class="flex items-center gap-3">
                            <img class="member-photo" :src="member.avatar_url" alt="">
                            <div class="member-declaration">
                                <span class="text-[hsl(0,0%,50%)] tracking-[-0.4px]">Декларация</span>
                                <span>до 26.12.2025</span>
                            </div>
                        </div>
                        <div class="self-start text-2xl"><ArrowIcon/></div>
                    </div>
                    <div class="member-about">
                        <div class="member-name">{{ member.fio }}</div>
                        <div class="member-income">
                            <span>{{ member.annual_income }}М ₽ / год</span>
                            <div class="dot"></div>
                            <span class="niche">{{ member.niche }}</span>
                        </div>
                    </div>
                </RouterLink>
            </div>
        </div>
    </section>
</template>

<style scoped>
.header {
    @apply text-[36px] font-medium tracking-[-1px] flex items-center p-3 gap-4
}

.members-container{
    @apply flex flex-col gap-1 justify-center align-top
}

.member-card{
    @apply flex flex-col bg-white rounded-[16px] p-[16px] gap-5 cursor-pointer
}

.member-card{
    font-family: 'SF Pro Text', Roboto, emoji, sans-serif;
}

.member-header{
    @apply flex flex-row items-center justify-between w-full tracking-[-0.4px]
}

.member-photo{
    @apply rounded-[50%] w-[48px] h-[48px]
}

.member-declaration{
    @apply flex flex-col text-[16px] leading-6
}

.members-about{
    @apply flex flex-col gap-1
}

.member-name {
    @apply text-xl tracking-[-0.4px] leading-6
}
.member-name {
    font-family: 'SF Pro Display', Roboto, emoji, sans-serif;
}
.member-income{
    @apply flex flex-row gap-1 text-[14px] font-[500] tracking-[-0.25px] leading-5
}
.dot {
    @apply w-[4px] h-[4px] bg-black inline-block rounded-[50%] self-center
}
</style>