<script setup lang="ts">
import type { Member } from '@/types/Member';
import { computed, onMounted,ref } from 'vue';
const props = defineProps<{
    members: Member[];
    current_user: Member
}>()

const displayedMembers = computed(() => props.members.slice(0, 9));
const extraCount = computed(() => props.members.length > 10 ? props.members.length - 9 : 0)

const areImagesLoaded = ref(false); 
const preloadImages = async () => {
  const imagePromises = displayedMembers.value.map((member) => {
    return new Promise((resolve, reject) => {
      const img = new Image();
      img.src = member.avatar_url;
      img.onload = resolve;
      img.onerror = reject;
    });
  });

  try {
    await Promise.all(imagePromises); 
    areImagesLoaded.value = true; 
  } catch (error) {
    console.error('Ошибка загрузки изображений:', error);
    areImagesLoaded.value = true; 
  }
};

// Запускаем загрузку при монтировании
onMounted(() => {
  preloadImages();
});
preloadImages();


</script>

<template>
    <div class="relative flex flex-wrap w-[100%]">
      <template v-for="(member, index) in displayedMembers" :key="member.username">
        <div class="avatar placeholder" :style="{ zIndex: 10 - index }">
          <img
            :src="member.avatar_url"
            alt="avatar"
            loading="eager"
            :class="[
              'avatar-img',
              current_user.username === member.username ? 'active-member' : '',
            ]"
          />
        </div>
      </template>
      <div v-if="extraCount > 0" class="extra-count avatar">+{{ extraCount }}</div>
    </div>
  </template>
<style scoped>
.avatar{
    @apply w-[38px] h-[38px] rounded-full border-2 border-white -mr-2 -mt-2
}
.placeholder {
  @apply bg-gray-300; 
}

.avatar-img {
  @apply w-full h-full rounded-full object-cover opacity-0 transition-opacity duration-300;
}
.avatar-img[src] {
  @apply opacity-100;
}
.active-member{
    @apply border-2 border-[#FF6644] z-10
}
.extra-count{
    @apply w-[38px] h-[38px] rounded-full bg-gray-300 text-base text-gray-700 font-medium flex items-center justify-center
}
</style>