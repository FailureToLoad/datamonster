import axios from "axios";
const ax = axios.create({
    headers: {
      "Accept-Language": "en-US,en;q=0.5",
    },
  });

export type Settlement = {
    name: string
    type: string
    limit: number
    departing: number
}

export type CreateSettlementRequest = {
    settlementName: string
    userId: string
}

type api = {
  getSettlementsForUser: (userId: string) => Promise<Settlement[]>
  getSettlement: () => Promise<Settlement>
  createSettlement: (request: CreateSettlementRequest) => Promise<Settlement>
    
}

const requester:api = {
  getSettlementsForUser: async function (userId: string): Promise<Settlement[]> {
      const response = await ax.get<Settlement[]>(`http://localhost:8000/settlement`,{headers: {Authorization: userId}})
      return response.data
  },
  getSettlement: async function (): Promise<Settlement> {
      const response = await ax.get<Settlement>('http://localhost:8000/settlement')
      return response.data
  },
  createSettlement: async function (request: CreateSettlementRequest): Promise<Settlement> {
    const response = await ax.post<Settlement>('http://localhost:8000/settlement', request)
    return response.data
  }
}

export default requester

