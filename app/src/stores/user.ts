import { defineStore } from "pinia"
import type { Member } from "@/types/Member"

export const useUserStore = defineStore("user", () => {
  const currentUser: Member = {
    fio: "Andrew Green",
    avatar_url: "/images/members/current_member.png",
    niche: "IT",
    annual_income: 100,
    username: "incetro",
  }

  return {
    currentUser,
  }
})
