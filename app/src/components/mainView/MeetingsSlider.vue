<script setup lang="ts">
import {Swiper, SwiperSlide} from 'swiper/vue';
import MeetingCard from './MeetingCard.vue'
import { Meeting } from '@/types/Meeting'

const props = defineProps<{
  meetings: Meeting[]
}>()
</script>

<template>
  <div class="slider-container">
    <Swiper
      :slides-per-view="1"
      :space-between="0"
      :initial-slide="3"
      :pagination="{ el: '.swiper-pagination', type: 'bullets', bulletClass: 'dot', bulletActiveClass: 'active' }"
      :centered-slides="true"
      class="swiper"
    >
      <SwiperSlide
        v-for="(meeting, index) in meetings"
        :key="index"
        class="slider-item"
      >
        <MeetingCard
          :meeting="meeting"
          :total="meetings.length"
          :currentIndex="index"
        />
      </SwiperSlide>
      <div class="swiper-pagination dots"></div>
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