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
    <div class="relative flex flex-wrap w-full">
      <template v-for="(member, index) in displayedMembers" :key="member.username">
        <img
          :src="member.avatar_url"
          alt="avatar"
          :class="[
            'avatar',
            current_user.username === member.username ? 'active-member' : '',
            index >= 5 ? 'underlap-row' : ''
          ]"
          :style="{ zIndex: 100 - index }"
        />
      </template>
      <div
        v-if="extraCount > 0"
        class="extra-count avatar underlap-row"
        :style="{ zIndex: 100 - displayedMembers.length }"
      >
        +{{ extraCount }}
      </div>
    </div>
  </template>
<style scoped>
.avatar {
  @apply w-9 h-9 rounded-full border-2 border-white -ml-2;
}

.active-member {
  @apply border-[#FF6644];
}

.extra-count {
  @apply bg-gray-300 text-sm text-gray-700 flex items-center justify-center;
}

.underlap-row {
  @apply -mt-2;
  position: relative;
}
</style>