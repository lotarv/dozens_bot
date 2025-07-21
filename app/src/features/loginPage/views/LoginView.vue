<script lang="ts" setup>

function openChatBot() {
    const botName = import.meta.env.VITE_BOT_NAME;
    const url = `https://t.me/${botName}`;

    if (window.Telegram && window.Telegram.WebApp) {
        Telegram.WebApp.openTelegramLink(url);
        Telegram.WebApp.close();
    } else {
        window.location.href = url;
    }
}

function openLink() {
    const url = `https://t.me/incetro`;
    if (window.Telegram && window.Telegram.WebApp) {
        Telegram.WebApp.openTelegramLink(url);
        Telegram.WebApp.close();
    } else {
        window.location.href = url;
    }
}
</script>
<template>
    <main class="login-main">
        <div class="info-block">
            <img src="../assets/logo.png" alt="" class="w-12 h-12">
            <p class="fond-medium text-[28px] leading-[32px] tracking-[-0.4px]">Упс</p>
            <p class="font-medium text-gray-400 text-base tracking-[-0.4px] text-center px-2">Пока что бот доступен только для некоторых Десяток. Если хотите попробовать, пишите <span class="font-bold" @click="openLink">@incetro</span></p>
        </div>
    </main>
</template>
<style scoped>
.login-main {
    @apply bg-center bg-cover flex flex-col items-center justify-center px-4 overflow-hidden rounded-xl;
    position:relative;
    z-index: 1;
}
.login-main::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-image: url(../assets/login-bg.png);
    background-size: cover;
    background-position: center;
    filter: blur(30px); /* Уровень размытия, можно настроить */
    z-index: -1; /* Помещаем фон под содержимое */
}


.info-block {
    @apply bg-white rounded-2xl flex flex-col gap-2 items-center text-center p-6
}

.go-button {
    @apply bg-[#d35e43] px-4 py-2 font-medium text-xl text-white rounded-xl;
}
</style>Fma