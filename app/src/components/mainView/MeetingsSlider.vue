<script setup lang="ts">
import { ref } from 'vue'
import MeetingCard from './MeetingCard.vue'
import { Meeting } from '@/types/Meeting'

const props = defineProps<{ meetings: Meeting[] }>()
const currentIndex = ref(0)
let touchStartX = 0
let touchEndX = 0

function onTouchStart(e: TouchEvent) {
  touchStartX = e.touches[0].clientX
}
function onTouchMove(e: TouchEvent) {
  touchEndX = e.touches[0].clientX
}
function onTouchEnd() {
  const delta = touchEndX - touchStartX
  if (delta > 50 && currentIndex.value > 0) {
    currentIndex.value--
  } else if (delta < -50 && currentIndex.value < props.meetings.length - 1) {
    currentIndex.value++
  }
}
</script>

<template>
  <div class="slider-container">
    <div
      class="slider-track"
      :style="{ transform: `translateX(-${currentIndex * 100}%)` }"
      @touchstart="onTouchStart"
      @touchmove="onTouchMove"
      @touchend="onTouchEnd"
    >
      <div v-for="(meeting, index) in meetings" :key="index" class="slider-item">
        <MeetingCard :meeting="meeting" :total="meetings.length" :currentIndex="index" />
      </div>
    </div>

    <div class="dots">
      <div v-for="(m, i) in meetings" :key="i" class="dot" :class="{ active: i === currentIndex }" />
    </div>
  </div>
</template>

<style scoped>
.slider-container {
  overflow: hidden;
  position: relative;
  width: 100%;
}
.slider-track {
  display: flex;
  transition: transform 0.3s ease;
}
.slider-item {
  flex: 0 0 100%;
}
.dots {
  position: absolute;
  top: 0;
  right: 3%;
  display: flex;
  gap: 4px;
  margin-top: 1rem;
}
.dot {
  width: 8px;
  height: 8px;
  background: rgba(255, 255, 255, 0.4);
  border-radius: 50%;
}
.dot.active {
  width: 20px;
  background: white;
}
</style>
