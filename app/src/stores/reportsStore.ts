import { getTelegramInitData } from "@/services/auth";
import axios from "axios";
import { defineStore } from "pinia";
import { ref } from "vue";
import { UserReports } from "@/types/UserReports";
export const useReportsStore = defineStore("reports", () => {
    const reports = ref<Record<string, UserReports>>({})
    const isLoading = ref<Record<string, boolean>>({})
    const error = ref<Record<string, string | null>>({})

    async function fetchUserReports(username: string) {
        if (reports.value[username]) return
        isLoading.value[username] = true
        error.value[username] = null
        try {
            const response = await axios.get<UserReports>(`${import.meta.env.VITE_API_URL}/reports/${username}`, {
                headers: {
                    'X-Telegram-Init-Data': getTelegramInitData(),
                },
            });
            reports.value[username] = {
                username: response.data.username,
                avatar_url: response.data.avatar_url,
                reports: response.data.reports,
            }
        } catch (err) {
            error.value[username] = 'Failed to load reports';
            console.error(`Failed to fetch reports for ${username}:`, err);
        } finally {
            isLoading.value[username] = false
        }
    }

    return {
        reports,
        isLoading,
        error,
        fetchUserReports
    }
})