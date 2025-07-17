<script lang="ts" setup>

function openChatBot() {
    const botName = import.meta.env.VITE_BOT_NAME
    window.location.href = `https://t.me/${botName}`
}

</script>
<template>
    <main class="login-main">
        <div class="info-block">
            <p class="fond-medium text-[28px] leading-[32px] tracking-[-0.4px]">Мы знакомы?</p>
            <p class="font-medium text-gray-400 text-base tracking-[-0.4px]">Похоже, вы не являетесь членом десятки. Вы можете присоединиться к своей десятке по коду или создать новую в нашем боте!</p>
            <button class="go-button" @click="openChatBot">Погнали!</button>
        </div>
    </main>
</template>
<style scoped>
.login-main {
    @apply bg-center bg-cover flex flex-col items-center justify-center px-4;
    background-image: url(../assets/login-bg.png);
}

.info-block {
    @apply bg-white rounded-2xl flex flex-col gap-2 items-center text-center p-6
}

.go-button {
    @apply bg-[#d35e43] px-4 py-2 font-medium text-xl text-white rounded-xl;
}
</style>Fma