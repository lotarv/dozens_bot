import { defineStore } from "pinia";
import bg1 from "@/assets/images/slider-backgrounds/bg1.png"
import bg2 from "@/assets/images/slider-backgrounds/bg2.png"
import bg3 from "@/assets/images/slider-backgrounds/bg3.png"
import { ref } from "vue";
import { Meeting2 } from "@/types/Meeting";
import { DozensTransport } from "@/repository/http";

export const UseMeetingsStore = defineStore('meetings', () => {
    const backgrounds = [bg1,bg2,bg3]
    const meetings = ref<Meeting2[]>([])

    async function fetchMeetings() {
        meetings.value = await DozensTransport.GetMeetings()
    }
    return {
        backgrounds,
        meetings,
        fetchMeetings
    }
})