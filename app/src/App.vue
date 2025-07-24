<script setup lang="ts">
import { onBeforeMount } from 'vue';
import MainView from './views/MainView.vue';
declare const Telegram: any

onBeforeMount(() => {
    const tg = Telegram.WebApp

    try {
        tg.expand()
        tg.isVerticalSwipesEnabled = false
        tg.disableVerticalSwipes()
        tg.ready()
    } catch(e) {
        alert(e)
    }
})
</script>
<template>
  <router-view>
  </router-view>
</template>


<style scoped>
</style>
