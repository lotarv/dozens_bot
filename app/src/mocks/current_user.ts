interface User {
    id: number;
    full_name: string;
    avatar_url: string | null;
    niche: string | null;
    annual_income: number;
    telegram_id: string;
}
import type { Member } from "@/types/Member";
export const current_user: Member = {
    fio: "Andrew Green",
    avatar_url:'/images/members/current_member.png',
    niche: "IT",
    annual_income: 100,
    username: "incetro"
}