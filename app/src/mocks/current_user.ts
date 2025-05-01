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
    avatar_url:'https://img.notionusercontent.com/s3/prod-files-secure%2F9a2e0635-b9d4-4178-a529-cf6b3bdce29d%2Fdbd2a2a7-da84-44fa-b7f2-e87b8833c38f%2Fphoto_2025-04-29_10.09.55.jpeg/size/w=2000?exp=1746175503&sig=5gMvvj1cV1K9R0-piF7TWE9eg_GO-Ml-A6U1n8CSh2A&id=1e027b05-0404-8016-b135-f7c62fc2c749&table=block&userId=e615d320-7569-4aee-9c73-66dcaee3f416',
    niche: "IT",
    annual_income: 100,
    username: "incetro"
}