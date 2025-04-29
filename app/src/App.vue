<script setup lang="ts">
</script>

<template>
  <router-view v-slot="{ Component }">
    <keep-alive include="MembersView">
      <component :is="Component" />
    </keep-alive>
  </router-view>
</template>


<style scoped>
</style>
