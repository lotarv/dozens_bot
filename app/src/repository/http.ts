import api from "@/services/api"
import { DeclarationDocument } from "@/features/otherMemberProfile/entities/DeclarationDocument"
import { PiggyBank } from "@/features/piggyBank/entities/PiggyBank"
import { Meeting2 } from "@/types/Meeting"

export class DozensTransport {
    static async GetDeclarations(username: string): Promise<DeclarationDocument[]> {
        try {
            const {data} = await api.get<DeclarationDocument[]>(`/declarations/${username}`)
            return data
        } catch(error) {
            console.error("failed to get declarations: ", error)
            return []
        }
    }
    static async GetDeclarationByID(id: string): Promise<DeclarationDocument | null> {
        try {
            const {data} = await api.get<DeclarationDocument>(`/declaration/${id}`)
            console.log("REPOSITORY: ", data)
            return data
        } catch(error) {
            console.error("failed to get declarations: ", error)
            return null
        }
    }
    static async GetPiggyBank(): Promise<PiggyBank | null> {
        try {
            const {data} = await api.get<PiggyBank>('/piggy-bank')
            return data
        } catch(e) {
            console.error("failed to get piggy bank: ", e)
            return null
        }
    }

    static async NewBankTransaction() {
        try {
            await api.post('/bot/new-transaction')
        } catch(error) {
            console.error("failed to open bot: ", error)
        }
    }

    static async GetMeetings(): Promise<Meeting2[]> {
        try {
            const {data} = await api.get<Meeting2[]>("/meetings")
            return data
        } catch(e) {
            console.error("failed to get meetings: ", e)
            return []
        }
    }
}