import { signOutUser } from "@/auth/firebase";
import axios from "axios";
const ax = axios.create({
    headers: {
      "Accept-Language": "en-US,en;q=0.5",
    },
  });
ax.defaults.withCredentials = true;
ax.interceptors.response.use(
    (response) => {
      return response;
    },
    (error) => {
      if (error.response.status === 401) {
        signOutUser();
      }
      return Promise.reject(error);
    }
)

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
  authorize: (token:string) => Promise<number>
  getSettlementsForUser: () => Promise<Array<Settlement>>
  getSettlement: () => Promise<Settlement>
  createSettlement: (request: CreateSettlementRequest) => Promise<Settlement>
}

const requester: api = {
  authorize: async function (token: string): Promise<number> {
    const response = await ax.post('http://localhost:8080/authorize', {token: token})
    return response.status
  },
  getSettlementsForUser: async function (): Promise<Settlement[]> {
      const response = await ax.get<AllSettlementsResponse>(`http://localhost:8080/settlement`)
      return response.data.settlements
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

