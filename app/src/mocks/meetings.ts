import { Member } from "@/types/Member"
import { Meeting } from "@/types/Meeting"
const meetingsRev: Meeting[] = [
    {
        date: "2025-05-25T13:00:00Z",
        location: "The Kitchen",
        members: [
            {
                fio:"",
                niche: "",
                annual_income: 100,
                username: "member1",
                avatar_url: "/images/members/member1.png"
            },
            {
                fio:"",
                niche: "",
                annual_income: 100,
                username: "member2",
                avatar_url: "/images/members/member2.png",
            },
            {
                fio:"",
                niche: "",
                annual_income: 100,
                username: "member3",
                avatar_url: "/images/members/member3.png",
            },
        ]
    },
    {
        date: "2025-04-20T13:00:00Z",
        location: "Хорошие люди",
        members: [
            {
                fio:"",
                niche: "",
                annual_income: 100,
                username: "member1",
                avatar_url: "/images/members/member1.png"
            },
            {
                fio:"",
                niche: "",
                annual_income: 100,
                username: "member2",
                avatar_url: "/images/members/member2.png",
            },
            {
                fio:"",
                niche: "",
                annual_income: 100,
                username: "member3",
                avatar_url: "/images/members/member3.png",
            },
        ]
    },
    {
        date: "2025-03-15T13:00:00Z",
        location: "Regions",
        members: [
            {
                fio:"",
                niche: "",
                annual_income: 100,
                username: "member1",
                avatar_url: "/images/members/member1.png"
            },
            {
                fio:"",
                niche: "",
                annual_income: 100,
                username: "member2",
                avatar_url: "/images/members/member2.png",
            },
            {
                fio:"",
                niche: "",
                annual_income: 100,
                username: "member3",
                avatar_url: "/images/members/member3.png",
            },
        ]
    },
    {
        date: "2025-02-25T13:00:00Z",
        location: "Good place",
        members: [
            {
                fio:"",
                niche: "",
                annual_income: 100,
                username: "member1",
                avatar_url: "/images/members/member1.png"
            },
            {
                fio:"",
                niche: "",
                annual_income: 100,
                username: "member2",
                avatar_url: "/images/members/member2.png",
            },
            {
                fio:"",
                niche: "",
                annual_income: 100,
                username: "member3",
                avatar_url: "/images/members/member3.png",
            },
        ]
    },
]

export const meetings = meetingsRev.reverse()