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
                avatar_url: "https://img.notionusercontent.com/s3/prod-files-secure%2F9a2e0635-b9d4-4178-a529-cf6b3bdce29d%2F0c493f69-7f75-4bab-8fb9-696cbb0a3e9c%2Fmidas-hofstra-tidSLv-UaNs-unsplash.jpg/size/w=2000?exp=1746276583&sig=dtwG4hLErNIbAdYNKU-983gXifudwxw0ZyLkBnAA4AA&id=1e427b05-0404-801e-b5ca-c31f9634954d&table=block&userId=e615d320-7569-4aee-9c73-66dcaee3f416"
            },
            {
                fio:"",
                niche: "",
                annual_income: 100,
                username: "member2",
                avatar_url: "https://img.notionusercontent.com/s3/prod-files-secure%2F9a2e0635-b9d4-4178-a529-cf6b3bdce29d%2Fe81977ef-8103-446d-9041-b5aa4fae82f5%2Falexander-hipp-iEEBWgY_6lA-unsplash.jpg/size/w=2000?exp=1746276606&sig=GcjkRBhZX4p9_XRHWNBUbYK6NyoMI7vOsc_kEAIWG5U&id=1e427b05-0404-80a0-a9ba-df263b3992c6&table=block&userId=e615d320-7569-4aee-9c73-66dcaee3f416"
            },
            {
                fio:"",
                niche: "",
                annual_income: 100,
                username: "member3",
                avatar_url: "https://img.notionusercontent.com/s3/prod-files-secure%2F9a2e0635-b9d4-4178-a529-cf6b3bdce29d%2Fdbd2a2a7-da84-44fa-b7f2-e87b8833c38f%2Fphoto_2025-04-29_10.09.55.jpeg/size/w=2000?exp=1746275932&sig=X2t-YJMgewxvj8_9ucUnkAybwLgIVa-hkcZ-beb4CW8&id=1e027b05-0404-8016-b135-f7c62fc2c749&table=block&userId=e615d320-7569-4aee-9c73-66dcaee3f416"
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
                avatar_url: "https://img.notionusercontent.com/s3/prod-files-secure%2F9a2e0635-b9d4-4178-a529-cf6b3bdce29d%2F0c493f69-7f75-4bab-8fb9-696cbb0a3e9c%2Fmidas-hofstra-tidSLv-UaNs-unsplash.jpg/size/w=2000?exp=1746276583&sig=dtwG4hLErNIbAdYNKU-983gXifudwxw0ZyLkBnAA4AA&id=1e427b05-0404-801e-b5ca-c31f9634954d&table=block&userId=e615d320-7569-4aee-9c73-66dcaee3f416"
            },
            {
                fio:"",
                niche: "",
                annual_income: 100,
                username: "member2",
                avatar_url: "https://img.notionusercontent.com/s3/prod-files-secure%2F9a2e0635-b9d4-4178-a529-cf6b3bdce29d%2Fe81977ef-8103-446d-9041-b5aa4fae82f5%2Falexander-hipp-iEEBWgY_6lA-unsplash.jpg/size/w=2000?exp=1746276606&sig=GcjkRBhZX4p9_XRHWNBUbYK6NyoMI7vOsc_kEAIWG5U&id=1e427b05-0404-80a0-a9ba-df263b3992c6&table=block&userId=e615d320-7569-4aee-9c73-66dcaee3f416"
            },
            {
                fio:"",
                niche: "",
                annual_income: 100,
                username: "member3",
                avatar_url: "https://img.notionusercontent.com/s3/prod-files-secure%2F9a2e0635-b9d4-4178-a529-cf6b3bdce29d%2Fdbd2a2a7-da84-44fa-b7f2-e87b8833c38f%2Fphoto_2025-04-29_10.09.55.jpeg/size/w=2000?exp=1746275932&sig=X2t-YJMgewxvj8_9ucUnkAybwLgIVa-hkcZ-beb4CW8&id=1e027b05-0404-8016-b135-f7c62fc2c749&table=block&userId=e615d320-7569-4aee-9c73-66dcaee3f416"
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
                avatar_url: "https://img.notionusercontent.com/s3/prod-files-secure%2F9a2e0635-b9d4-4178-a529-cf6b3bdce29d%2F0c493f69-7f75-4bab-8fb9-696cbb0a3e9c%2Fmidas-hofstra-tidSLv-UaNs-unsplash.jpg/size/w=2000?exp=1746276583&sig=dtwG4hLErNIbAdYNKU-983gXifudwxw0ZyLkBnAA4AA&id=1e427b05-0404-801e-b5ca-c31f9634954d&table=block&userId=e615d320-7569-4aee-9c73-66dcaee3f416"
            },
            {
                fio:"",
                niche: "",
                annual_income: 100,
                username: "member2",
                avatar_url: "https://img.notionusercontent.com/s3/prod-files-secure%2F9a2e0635-b9d4-4178-a529-cf6b3bdce29d%2Fe81977ef-8103-446d-9041-b5aa4fae82f5%2Falexander-hipp-iEEBWgY_6lA-unsplash.jpg/size/w=2000?exp=1746276606&sig=GcjkRBhZX4p9_XRHWNBUbYK6NyoMI7vOsc_kEAIWG5U&id=1e427b05-0404-80a0-a9ba-df263b3992c6&table=block&userId=e615d320-7569-4aee-9c73-66dcaee3f416"
            },
            {
                fio:"",
                niche: "",
                annual_income: 100,
                username: "member3",
                avatar_url: "https://img.notionusercontent.com/s3/prod-files-secure%2F9a2e0635-b9d4-4178-a529-cf6b3bdce29d%2Fdbd2a2a7-da84-44fa-b7f2-e87b8833c38f%2Fphoto_2025-04-29_10.09.55.jpeg/size/w=2000?exp=1746275932&sig=X2t-YJMgewxvj8_9ucUnkAybwLgIVa-hkcZ-beb4CW8&id=1e027b05-0404-8016-b135-f7c62fc2c749&table=block&userId=e615d320-7569-4aee-9c73-66dcaee3f416"
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
                avatar_url: "https://img.notionusercontent.com/s3/prod-files-secure%2F9a2e0635-b9d4-4178-a529-cf6b3bdce29d%2F0c493f69-7f75-4bab-8fb9-696cbb0a3e9c%2Fmidas-hofstra-tidSLv-UaNs-unsplash.jpg/size/w=2000?exp=1746276583&sig=dtwG4hLErNIbAdYNKU-983gXifudwxw0ZyLkBnAA4AA&id=1e427b05-0404-801e-b5ca-c31f9634954d&table=block&userId=e615d320-7569-4aee-9c73-66dcaee3f416"
            },
            {
                fio:"",
                niche: "",
                annual_income: 100,
                username: "member2",
                avatar_url: "https://img.notionusercontent.com/s3/prod-files-secure%2F9a2e0635-b9d4-4178-a529-cf6b3bdce29d%2Fe81977ef-8103-446d-9041-b5aa4fae82f5%2Falexander-hipp-iEEBWgY_6lA-unsplash.jpg/size/w=2000?exp=1746276606&sig=GcjkRBhZX4p9_XRHWNBUbYK6NyoMI7vOsc_kEAIWG5U&id=1e427b05-0404-80a0-a9ba-df263b3992c6&table=block&userId=e615d320-7569-4aee-9c73-66dcaee3f416"
            },
            {
                fio:"",
                niche: "",
                annual_income: 100,
                username: "member3",
                avatar_url: "https://img.notionusercontent.com/s3/prod-files-secure%2F9a2e0635-b9d4-4178-a529-cf6b3bdce29d%2Fdbd2a2a7-da84-44fa-b7f2-e87b8833c38f%2Fphoto_2025-04-29_10.09.55.jpeg/size/w=2000?exp=1746275932&sig=X2t-YJMgewxvj8_9ucUnkAybwLgIVa-hkcZ-beb4CW8&id=1e027b05-0404-8016-b135-f7c62fc2c749&table=block&userId=e615d320-7569-4aee-9c73-66dcaee3f416"
            },
        ]
    },
]

export const meetings = meetingsRev.reverse()