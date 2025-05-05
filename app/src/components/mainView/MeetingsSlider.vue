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
let isTouching = false


function handlePointerDown(e: PointerEvent) {
  if (isTouching) return // если уже обрабатываем touch, игнорируем
  touchStartX = e.clientX
}

function handlePointerMove(e: PointerEvent) {
  if (isTouching) return
  touchEndX = e.clientX
}

function handlePointerUp() {
  if (isTouching) return
  handleSwipe()
}

function handleTouchStart(e: TouchEvent) {
  isTouching = true
  touchStartX = e.touches[0].clientX
}

function handleTouchMove(e: TouchEvent) {
  touchEndX = e.touches[0].clientX
}

function handleTouchEnd() {
  handleSwipe()
  setTimeout(() => {
    isTouching = false // снимаем флаг после завершения
  }, 0)
}

function handleSwipe() {
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
      @pointerdown="handlePointerDown"
      @pointermove="handlePointerMove"
      @pointerup="handlePointerUp"
      @touchstart="handleTouchStart"
      @touchmove="handleTouchMove"
      @touchend="handleTouchEnd"
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
  touch-action: manipulation;
  @apply relative w-full overflow-hidden;
}

.slider-track {
  touch-action: none !important;
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
