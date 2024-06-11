import api from "./api";

export type Settlement = {
  id: string;
  name: string;
  limit: number;
  departing: number;
  cc: number;
  year: number;
};

export type SettlementCreationRequest = {
  name: string;
};

type AllSettlementsResponse = {
  settlements: Array<Settlement>;
  count: number;
};

type settlementRequests = {
  getSettlementsForUser: () => Promise<Array<Settlement> | null>;
  getSettlement: (id: string) => Promise<Settlement | null>;
  createSettlement: (request: SettlementCreationRequest) => Promise<Settlement>;
};

const SettlementApi: settlementRequests = {
  getSettlementsForUser: async function (): Promise<Settlement[] | null> {
    try {
      const response = await api.get<AllSettlementsResponse | null>(
        `http://localhost:8080/settlement`,
      );
      if (!response.data) {
        return null;
      }
      return response.data.settlements;
    } catch (e) {
      console.log(e);
      return null;
    }
  },
  getSettlement: async function (id: string): Promise<Settlement | null> {
    try {
      const response = await api.get<Settlement | null>(
        "http://localhost:8080/settlement/" + id,
      );
      if (!response.data) {
        return null;
      }
      return response.data;
    } catch (e) {
      console.log(e);
      return null;
    }
  },
  createSettlement: async function (
    request: SettlementCreationRequest,
  ): Promise<Settlement> {
    const response = await api.post<Settlement>(
      "http://localhost:8080/settlement",
      request,
    );
    return response.data;
  },
};

export default SettlementApi;
