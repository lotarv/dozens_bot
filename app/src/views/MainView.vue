<script lang="ts" setup>
import axios from 'axios';
import { ref, onBeforeMount, Ref } from 'vue';

interface User {
    id: number;
    full_name: string;
    avatar_url: string | null;
    niche: string | null;
    annual_income: number;
    telegram_id: string;
}

const users: Ref<User[]> = ref([]); 

async function fetchUsers() {
    try {
        const response = await axios.get<User[]>(`${import.meta.env.VITE_API_URL}/users`);
        users.value = response.data; // Теперь TypeScript знает, что у users есть свойство .value
    } catch (err) {
        console.error('Failed to fetch users:', err);
    }
}

onBeforeMount(fetchUsers);
</script>

<template>
    <div>
        <div v-if="users.length > 0">
            <ul>
                <li v-for="user in users" :key="user.id">
                    {{ user.full_name }} (ID: {{ user.id }})
                </li>
            </ul>
        </div>
        <div v-else>No users in db</div>
    </div>
</template>

<style></style>