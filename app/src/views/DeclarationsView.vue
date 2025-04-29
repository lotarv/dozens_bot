<script lang="ts" setup>
import {useRoute} from 'vue-router'
const route = useRoute()
const username = route.params.username
</script>
<template>
    <section>
        <h1>Декларации пользователя: @{{ username }}</h1>
        <p>Раздел находится в разработке...</p>
        <RouterLink
        :to="'/'"
        class="bg-gray-300 rounded-lg p-2">Назад</RouterLink>
    </section>
</template>

<style scoped>
section {
  @apply flex flex-col justify-center items-center h-screen text-center;
}
h1 {
    @apply text-center tracking-[-0.4px] text-2xl
}
p{
    @apply text-center tracking-[-0.4px] text-2xl

}
</style>