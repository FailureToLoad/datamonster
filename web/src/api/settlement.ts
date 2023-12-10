import api from "./api";

export type Settlement = {
  id: string;
  name: string;
  limit: number;
  departing: number;
  cc: number;
  year: number;
};

export type CreateSettlementRequest = {
  name: string;
};

type AllSettlementsResponse = {
  settlements: Array<Settlement>;
  count: number;
};

type settlementApi = {
  getSettlementsForUser: () => Promise<Array<Settlement>>;
  getSettlement: (id: string) => Promise<Settlement>;
  createSettlement: (request: CreateSettlementRequest) => Promise<Settlement>;
};

const settlementRequests: settlementApi = {
  getSettlementsForUser: async function (): Promise<Settlement[]> {
    const response = await api.get<AllSettlementsResponse>(
      `http://localhost:8080/settlement`,
    );
    return response.data.settlements;
  },
  getSettlement: async function (id: string): Promise<Settlement> {
    const response = await api.get<Settlement>(
      "http://localhost:8080/settlement/" + id,
    );
    return response.data;
  },
  createSettlement: async function (
    request: CreateSettlementRequest,
  ): Promise<Settlement> {
    const response = await api.post<Settlement>(
      "http://localhost:8080/settlement",
      request,
    );
    return response.data;
  },
};

export default settlementRequests;
