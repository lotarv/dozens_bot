import { defineStore } from "pinia";
import bg1 from "@/assets/images/slider-backgrounds/bg1.png"
import bg2 from "@/assets/images/slider-backgrounds/bg2.png"
import bg3 from "@/assets/images/slider-backgrounds/bg3.png"

export const UseMeetingsStore = defineStore('meetings', () => {
    const backgrounds = [bg1,bg2,bg3]
    return {backgrounds}
})