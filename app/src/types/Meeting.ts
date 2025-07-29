import type { Member } from "./Member"
export interface Meeting{
    date: string,
    location: string,
    members: Member[],
    map_url:string
}

export interface Meeting2 {
    start_time: string,
    location_name:string,
    map_url: string
}