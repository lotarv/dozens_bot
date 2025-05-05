<script setup lang="ts">
import { ref } from 'vue'
import MeetingCard from './MeetingCard.vue'
import { Meeting } from '@/types/Meeting'

const props = defineProps<{
  meetings: Meeting[]
}>()

const currentIndex = ref(3)

let touchStartX = 0
let touchEndX = 0

function handlePointerDown(e: PointerEvent) {
  touchStartX = e.clientX;
}

function handlePointerMove(e: PointerEvent) {
  touchEndX = e.clientX;
}

function handlePointerUp() {
  const delta = touchEndX - touchStartX;
  if (delta > 50 && currentIndex.value > 0) {
    currentIndex.value--;
  } else if (delta < -50 && currentIndex.value < props.meetings.length - 1) {
    currentIndex.value++;
  }
}
</script>
<template>
  <div class="slider-container">
    <div
      class="slider-track"
      :style="{ transform: `translateX(-${currentIndex * 100}%)` }"
      @pointerdown="handlePointerDown"
      @pointermove="handlePointerMove"
      @pointerup="handlePointerUp"
    >
      <div
        v-for="(meeting, index) in meetings"
        :key="index"
        class="slider-item"
      >
        <MeetingCard
          :meeting="meeting"
          :total="meetings.length"
          :currentIndex="index"
        />
      </div>
    </div>

    <div class="dots">
      <div
        v-for="(m, i) in meetings"
        :key="i"
        class="dot"
        :class="{ active: i === currentIndex }"
      ></div>
    </div>
  </div>
</template>

<style scoped>
.slider-container {
  touch-action: pan-y;
  @apply relative w-full overflow-hidden;
}

.slider-track {
  will-change: transform;
  @apply flex transition-transform duration-500 ease-in-out;
  width: 100%;
}

.slider-item {
  flex: 0 0 100%;
  @apply box-border;
}

.dots {
  @apply flex justify-center gap-1 mt-3 absolute top-0 right-3;
}

.dot {
  @apply w-2 h-2 bg-white/40 rounded-full transition-all;
}

.dot.active {
  @apply w-5 bg-white;
}
</style>
