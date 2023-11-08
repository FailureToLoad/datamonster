import axios from "axios";
const ax = axios.create({
    headers: {
      "Accept-Language": "en-US,en;q=0.5",
    },
  });

export type Settlement = {
    id: string
    name: string
    limit: number
    departing: number
    cc: number
    year: number
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
  getSettlementsForUser: async function (token: string): Promise<Settlement[]> {
      const response = await ax.get<Settlement[]>(`http://localhost:8080/settlement`,{headers: {'Authorization': `${token}`}})
      return response.data
  },
  getSettlement: async function (): Promise<Settlement> {
      const response = await ax.get<Settlement>('http://localhost:8080/settlement')
      return response.data
  },
  createSettlement: async function (request: CreateSettlementRequest): Promise<Settlement> {
    const response = await ax.post<Settlement>('http://localhost:8080/settlement', request)
    return response.data
  }
}

export default requester

