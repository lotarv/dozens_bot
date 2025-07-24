import { getTelegramInitData } from "@/services/auth";
import axios from "axios";
import { defineStore } from "pinia";
import { ref } from "vue";
import { UserReports } from "@/types/UserReports";
import { decryptGoAES, isEncrypted } from "@/services/crypto";
import { useDecryptionStore } from "./decryption";
export const useReportsStore = defineStore("reports", () => {
    const cryptoStore = useDecryptionStore()
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
//     async function fetchUserReports(username: string) {
//     if (reports.value[username]) return

//     isLoading.value[username] = true
//     error.value[username] = null

//     try {
//         const response = await axios.get<UserReports>(`${import.meta.env.VITE_API_URL}/reports/${username}`, {
//             headers: {
//                 'X-Telegram-Init-Data': getTelegramInitData(),
//             },
//         })

//         const rawReports = Array.isArray(response.data.reports) ? response.data.reports : []
        
//         // Группируем по дате
//         const groupedReports: Record<string, string[]> = {}

//         for (const report of rawReports) {
//             const date = report.creation_date
//             if (!groupedReports[date]) groupedReports[date] = []

//             let text: string
//             if (isEncrypted(report.text)) {
//                 // Расшифровка
//                 text = await decryptGoAES(report.text, cryptoStore.key)
//             } else {
//                 text = report.text
//             }

//             groupedReports[date].push(text)
//         }

//         // Сортируем даты и объединяем тексты в один отчёт на дату
//         const mergedReports = Object.entries(groupedReports)
//             .sort(([dateA], [dateB]) => new Date(dateB).getTime() - new Date(dateA).getTime()) // если нужен порядок по убыванию
//             .map(([creation_date, texts]) => ({
//                 creation_date,
//                 text: texts.join('\n\n'), // можно '\n---\n' или с другим разделителем
//             }))

//         reports.value[username] = {
//             username: response.data.username,
//             avatar_url: response.data.avatar_url,
//             reports: mergedReports,
//         }
//     } catch (err) {
//         error.value[username] = 'Failed to load reports'
//         console.error(`Failed to fetch reports for ${username}:`, err)
//     } finally {
//         isLoading.value[username] = false
//     }
// }



    return {
        reports,
        isLoading,
        error,
        fetchUserReports
    }
})