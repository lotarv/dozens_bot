import { defineStore } from "pinia";
import { ref } from "vue";
import { PiggyBank } from "../entities/PiggyBank";
import { DozensTransport } from "@/repository/http";

export const useBankStore = defineStore('piggy-bank', () => {
    const bank = ref<PiggyBank | null> (null)

    const fetchPiggyBank = async() => {
        bank.value = await DozensTransport.GetPiggyBank()
    }

    return {
        bank,
        fetchPiggyBank
    }
})