import type { Member } from "./Member"
export interface Meeting{
    date: string,
    location: string,
    members: Member[],
    map_url:string
}