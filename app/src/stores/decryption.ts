import api from "@/services/api";
import { defineStore } from "pinia";
import { ref } from "vue";

export const useDecryptionStore = defineStore('decryption', () => {
    const key = ref<string>("")
    
    async function fetchKey() {
        if (key.value) return
        try {
            const {data} = await api.get("/users/enc-key")
            key.value = data.key
            console.log("GOT ENCRYPTION KEY: ", key.value)
        } catch(e) {
            console.error("failed to fetch decryption key: ", e)
        }
    }
    return {
        key,
        fetchKey
    }
})