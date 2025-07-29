<script setup lang="ts">
import {Swiper, SwiperSlide} from 'swiper/vue';
import MeetingCard from './MeetingCard.vue'
import { Meeting, Meeting2 } from '@/types/Meeting'
import { onMounted, ref } from 'vue';

const props = defineProps<{
  meetings: Meeting2[]
}>()

const currentIndex = ref(props.meetings.length - 1)

function handleSlideChange(swiper:any) {
  currentIndex.value = swiper.realIndex
}

</script>


<template>
  <div class="slider-container">
    <Swiper
      :slides-per-view="1"
      :space-between="4"
      :initial-slide="currentIndex"
      :centered-slides="true"
      class="swiper"
      @slide-change="handleSlideChange"
    >
  <SwiperSlide v-for="(meeting, index) in meetings" :key="index" class="slider-item">
    <MeetingCard :meeting="meeting" :total="meetings.length" :currentIndex="index" />
  </SwiperSlide>

  <!-- Пагинация -->
  <div class="swiper-pagination dots">
    <div
        v-for="(m, i) in meetings"
        :key="i"
        class="dot"
        :class="{ active: i === currentIndex }"
      ></div>
  </div>
</Swiper>
  </div>
</template>

<style scoped>
.slider-container {
  @apply relative w-full overflow-hidden;
}

.swiper {
  @apply w-full;
}

.slider-item {
  @apply box-border w-full;
}

.dots {
  @apply flex justify-center gap-1 mt-3 absolute top-0 right-3;
}

.dot {
  @apply w-2 h-2 bg-white/40 rounded-full transition-all cursor-pointer;
}

.dot.active {
  @apply w-5 bg-white;
}
</style>