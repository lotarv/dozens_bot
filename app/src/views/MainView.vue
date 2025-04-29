<script lang="ts" setup>
import axios from 'axios';
import { ref, onBeforeMount, Ref } from 'vue';
import { getTelegramInitData } from '../services/auth';

interface User {
    id: number;
    full_name: string;
    avatar_url: string | null;
    niche: string | null;
    annual_income: number;
    telegram_id: string;
}

const users: Ref<User[]> = ref([]);
const isLoading = ref(false);
const error = ref<string | null>(null);

async function createOrUpdateUser() {
    try {
        await axios.post(
            `${import.meta.env.VITE_API_URL}/users`,
            {},
            {
                headers: {
                    'X-Telegram-Init-Data': getTelegramInitData(),
                },
            }
        );
    } catch (err) {
        error.value = 'Failed to authenticate. Please try again.';
        console.error('Authentication failed:', err);
    }
}

async function fetchUsers() {
    try {
        isLoading.value = true;
        error.value = null;

        const response = await axios.get<User[]>(`${import.meta.env.VITE_API_URL}/users`, {
            headers: {
                'X-Telegram-Init-Data': getTelegramInitData(),
            },
        });
        users.value = response.data;
        console.log(users)
    } catch (err) {
        error.value = 'Failed to load users. Please try again later.';
        console.error('Failed to fetch users:', err);
    } finally {
        isLoading.value = false;
    }
}

async function fetchMembers() {
    try {
        isLoading.value = true;
        error.value = null;

        const response = await axios.get<User[]>(`${import.meta.env.VITE_API_URL}/members`, {
            headers: {
                'X-Telegram-Init-Data': getTelegramInitData(),
            },
        });
        users.value = response.data;
        console.log(users)
    } catch (err) {
        error.value = 'Failed to load users. Please try again later.';
        console.error('Failed to fetch users:', err);
    } finally {
        isLoading.value = false;
    }
}

onBeforeMount(async () => {
    await createOrUpdateUser();
    await fetchUsers();
});
</script>

<template>
    <h1>Dozen's bot</h1>
    <p>В процессе разработки...</p>
</template>

<style>
.error {
    color: red;
}
ul {
    padding:0;
    margin:0;
}
li{
    padding:0;
    margin:0;
}
</style>