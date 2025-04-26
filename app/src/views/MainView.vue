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

onBeforeMount(async () => {
    await createOrUpdateUser();
    await fetchUsers();
});
</script>

<template>
    <div>
        <div v-if="isLoading">Loading...</div>
        <div v-else-if="error" class="error">{{ error }}</div>
        <div v-else-if="users.length > 0">
            <ul>
                <li v-for="user in users" :key="user.id">username (ID: {{ user.id }})
                </li>
            </ul>
        </div>
        <div v-else>No users in db</div>
    </div>
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