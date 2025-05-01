<script setup lang="ts">
import type { Member } from '@/types/Member';
import { computed } from 'vue';
const props = defineProps<{
    members: Member[];
    current_user: Member
}>()

const displayedMembers = computed(() => props.members.slice(0, 9));
const extraCount = computed(() => props.members.length > 10 ? props.members.length - 9 : 0)
</script>

<template>
    <div class="relative flex flex-wrap w-[100%]">
        <template v-for="(member, index) in displayedMembers" :key="member.username">
            <img 
            :src="member.avatar_url" alt="avatar"
            :class="[
            'avatar',
            current_user.username === member.username ? 'active-member' : '',
          ]"
          :style="{zIndex: 10 - index}">
        </template>
        <div
        v-if="extraCount > 0"
        class="extra-count avatar"> + {{ extraCount }}</div>
    </div>
</template>

<style scoped>
.avatar{
    @apply w-[38px] h-[38px] rounded-full border-2 border-white -mr-2 -mt-2
}
/* .avatar:first-child{
    @apply ml-0
} */
.active-member{
    @apply border-2 border-[#FF6644] z-10
}
.extra-count{
    @apply w-[38px] h-[38px] rounded-full bg-gray-300 text-sm text-gray-700 flex items-center justify-center
}
</style>