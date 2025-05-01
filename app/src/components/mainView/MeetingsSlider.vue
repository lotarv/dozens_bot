<template>
    <div class="meetings-slider">
      <Carousel
        v-bind="carouselConfig"
        @slide-end="changeBackground"
        :style="{ backgroundImage: `url(${currentBackground})` }"
      >
        <Slide v-for="(meeting, index) in meetings" :key="index">
          <div class="meeting-card">
            <h2 class="meeting-date">{{ formatDate(meeting.date) }}</h2>
            <h3 class="meeting-location">{{ meeting.location }}</h3>
            <div class="members-list">
              <div v-for="member in meeting.members" :key="member.username" class="member">
                <img :src="member.avatar_url" :alt="member.username" class="member-avatar" />
                <span class="member-name">{{ member.username }}</span>
              </div>
            </div>
          </div>
        </Slide>
        <template #addons>
          <Navigation />
          <Pagination />
        </template>
      </Carousel>
    </div>
  </template>
  
  <script>
  import { meetings } from '@/mocks/meetings'
  import 'vue3-carousel/dist/carousel.css'
  import { Carousel, Slide, Pagination, Navigation } from 'vue3-carousel'
  
  export default {
    name: 'MeetingsSlider',
    components: {
      Carousel,
      Slide,
      Pagination,
      Navigation,
    },
    data() {
      return {
        meetings,
        currentBackground: '',
        backgrounds: [
          '../../assets/images/slider-backgrounds/background1',
          '../../assets/images/slider-backgrounds/background2',
          '../../assets/images/slider-backgrounds/background3',
        ],
        carouselConfig: {
          itemsToShow: 1,
          snapAlign: 'center',
          wrapAround: false,
          transition: 500, // Плавная анимация в миллисекундах
        },
      }
    },
    methods: {
      formatDate(dateString) {
        const date = new Date(dateString)
        return date.toLocaleDateString('en-US', {
          year: 'numeric',
          month: 'long',
          day: 'numeric',
          hour: '2-digit',
          minute: '2-digit',
        })
      },
      changeBackground() {
        const randomIndex = Math.floor(Math.random() * this.backgrounds.length)
        this.currentBackground = this.backgrounds[randomIndex]
      },
    },
    mounted() {
      // Установить начальный фон
      this.changeBackground()
    },
  }
  </script>
  
  <style scoped>
  .meetings-slider {
    max-width: 600px;
    margin: 0 auto;
    padding: 20px;
    background-size: cover;
    background-position: center;
    transition: background-image 0.5s ease-in-out;
    position: relative;
    min-height: 400px;
  }
  
  .meeting-card {
    background: rgba(245, 245, 245, 0.9); /* Полупрозрачный фон для читаемости */
    border-radius: 8px;
    padding: 20px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    text-align: center;
    width: 100%;
  }
  
  .meeting-date {
    font-size: 1.5rem;
    color: #333;
    margin-bottom: 10px;
  }
  
  .meeting-location {
    font-size: 1.2rem;
    color: #666;
    margin-bottom: 20px;
  }
  
  .members-list {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: 20px;
  }
  
  .member {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  .member-avatar {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    object-fit: cover;
    margin-bottom: 8px;
  }
  
  .member-name {
    font-size: 0.9rem;
    color: #333;
  }
  
  /* Стили для навигации и пагинации vue3-carousel */
  .carousel__prev,
  .carousel__next {
    background: #007bff;
    color: white;
    border: none;
    border-radius: 50%;
    width: 40px;
    height: 40px;
    font-size: 1.2rem;
    transition: background 0.3s;
  }
  
  .carousel__prev:hover,
  .carousel__next:hover {
    background: #0056b3;
  }
  
  .carousel__prev--disabled,
  .carousel__next--disabled {
    background: #ccc;
    cursor: not-allowed;
  }
  
  .carousel__pagination {
    display: flex;
    justify-content: center;
    gap: 10px;
    margin-top: 20px;
  }
  
  .carousel__pagination-button {
    width: 10px;
    height: 10px;
    background: #ccc;
    border-radius: 50%;
    cursor: pointer;
  }
  
  .carousel__pagination-button--active {
    background: #007bff;
  }
  </style>