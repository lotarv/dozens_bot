import api from "@/services/api";
import { defineStore } from "pinia";
import { ref } from "vue";

export const useDecryptionStore = defineStore('decryption', () => {
    const pepper = import.meta.env.VITE_ENCRYPTION_PEPPER
    const key = ref<string>("")
    
    async function fetchKey() {
        if (key.value) return
        try {
            const {data} = await api.get("/users/dozen-code")
            key.value = `${data.code}${pepper}`
        } catch(e) {
            console.error("failed to fetch decryption key: ", e)
        }
    }
    return {
        key,
        fetchKey
    }
})