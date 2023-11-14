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
    name: string
}

type AllSettlementsResponse = {
    settlements: Array<Settlement>
    count: number
}

type api = {
  getSettlementsForUser: (userId: string) => Promise<Array<Settlement>>
  getSettlement: () => Promise<Settlement>
  createSettlement: (requesr: CreateSettlementRequest, token: string) => Promise<Settlement>
    
}

const requester: api = {
  getSettlementsForUser: async function (token: string): Promise<Settlement[]> {
      const response = await ax.get<AllSettlementsResponse>(`http://localhost:8080/settlement`,{headers: {'Authorization': `${token}`}})
      return response.data.settlements
  },
  getSettlement: async function (): Promise<Settlement> {
      const response = await ax.get<Settlement>('http://localhost:8080/settlement')
      return response.data
  },
  createSettlement: async function (request: CreateSettlementRequest, token: string): Promise<Settlement> {
    const response = await ax.post<Settlement>('http://localhost:8080/settlement', request, {headers: {'Authorization': `${token}`}})
    return response.data
  }
}

export default requester

