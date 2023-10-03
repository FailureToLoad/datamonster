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

type api = {
  getSettlement: () => Promise<Settlement>
    
}

const requester:api = {
  getSettlement: async function (): Promise<Settlement> {
      const response = await ax.get<Settlement>('http://localhost:8000/settlement')
      return response.data
  }
}

export default requester

