import api from "@/services/api"
import { DeclarationDocument } from "@/features/otherMemberProfile/entities/DeclarationDocument"

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
}