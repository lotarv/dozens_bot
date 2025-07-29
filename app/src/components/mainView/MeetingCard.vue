<script setup lang="ts">
import { Meeting, Meeting2 } from '@/types/Meeting';
import ArrowIcon from '../icons/ArrowIcon.vue';
import { computed } from 'vue';
import { UseMeetingsStore } from '@/stores/meetingsStore';
import { useMembersStore } from '@/stores/membersStore';

const meetingsStore = UseMeetingsStore()
const props = defineProps<{
    meeting: Meeting2;
    currentIndex: number,
    total: number,
}>()


const colors = [
    "#FF6644",
    "#657CFF",
    "#58ca67"
]
function formatDateTime(isoDate: string): string {
    const date = new Date(isoDate);
    
    // Массив месяцев на русском
    const months = [
        "января", "февраля", "марта", "апреля", "мая", "июня",
        "июля", "августа", "сентября", "октября", "ноября", "декабря"
    ];
    
    // Получаем день, месяц и время
    const day = date.getUTCDate();
    const month = months[date.getUTCMonth()];
    const hours = date.getUTCHours().toString().padStart(2, "0");
    const minutes = date.getUTCMinutes().toString().padStart(2, "0");
    
    return `${day} ${month}, ${hours}:${minutes}`;
}

const isNextMeeting = computed(() => {
    const meetingDate = new Date(props.meeting.start_time)
    const now = new Date();
    return meetingDate > now;
})

function openLink(url: string) {
    if (url) {
        window.open(url, '_blank');
    }
}

const members =  [
    {
        fio:"",
        niche: "",
        annual_income: 100,
        username: "member1",
        avatar_url: "/images/members/member1.png"
    },
    {
        fio:"",
        niche: "",
        annual_income: 100,
        username: "member2",
        avatar_url: "/images/members/member2.png",
    },
    {
        fio:"",
        niche: "",
        annual_income: 100,
        username: "member3",
        avatar_url: "/images/members/member3.png",
    },
]
</script>
<template>
    <div class="card" :style="{
        backgroundImage: `url(${meetingsStore.backgrounds[currentIndex % meetingsStore.backgrounds.length]})`
    }">
        <div class="card-header">
            <div class="meeting-status">
                <p v-if="isNextMeeting">Следующая встреча</p>
                <p v-else>Предыдущая встреча</p>
            </div>
        </div>
        <div class="content">
            <div class="avatars">
                <img
                    v-for="member in members"
                    :src="member.avatar_url"
                    alt="avatar"
                    :style="{borderColor: colors[currentIndex % colors.length]}">
            </div>
            <div class="meeting-date">{{ formatDateTime(meeting.start_time) }}</div>
            <div class="meeting-location">
                <span @click="openLink(meeting.map_url)">{{ meeting.location_name }}</span>
                <span class="text-[24px]"><ArrowIcon></ArrowIcon></span>
            </div>
        </div>
    </div>
</template>

<style scoped>

.card{
    @apply flex flex-col w-full justify-between p-4 gap-8 rounded-[16px] tracking-[-0.4px] bg-cover bg-center bg-no-repeat
}
.card-header{
    font-family: "SF Pro Text";
    @apply text-white leading-6 tracking-[-0.4px] font-medium text-base
}
.avatars{
    @apply flex mb-2
}
.avatars img {
    @apply h-12 w-12 rounded-full -mr-4 border-4
}
.card-header{
    @apply flex justify-between
}


.active {
    @apply w-4 bg-white
}
.content {
    font-family: "SF Pro Display", Arial, sans-serif;
    @apply flex flex-col font-medium text-[28px] leading-8
}

.meeting-date{
    @apply text-white
}

.meeting-location {
    @apply flex gap-2 items-center
}
</style>