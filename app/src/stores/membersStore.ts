import {defineStore} from 'pinia'
import {ref} from 'vue'
import axios from 'axios'
import { getTelegramInitData } from '@/services/auth'

import type { Member } from '@/types/Member'

export const useMembersStore = defineStore('members', () => {
    const members = ref<Member[]>([])
    const isLoading = ref(false)
    const error = ref<string | null>(null)

    async function fetchMembers() {
        if (members.value.length > 0) return
        try {
            isLoading.value = true;
            error.value = null;
            console.log(`making a request for ${import.meta.env.VITE_API_URL}/members`)
    
            const response = await axios.get<Member[]>(`${import.meta.env.VITE_API_URL}/members`, {
                headers: {
                    'X-Telegram-Init-Data': getTelegramInitData(),
                },
            });
            members.value = response.data.reverse();
            console.log(members)
        } catch (err) {
            error.value = 'Failed to load members. Please try again later.';
            console.error('Failed to fetch members:', err);
        } finally {
            isLoading.value = false
        }
    }

    function getMemberByUsername(username:string): Member{
        for (let member of members.value) {
            if (member.username == username) {
                return member
            }
        }

        return {
            username: '',
            fio: '',
            avatar_url: '',
            annual_income: 0,
            niche: ''
        } as Member
    }

    return {members, isLoading, error, fetchMembers, getMemberByUsername}
})