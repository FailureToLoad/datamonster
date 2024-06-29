import axios from "axios";

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
      const response = await axios.get<AllSettlementsResponse | null>(
        `settlement`,
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
      const response = await axios.get<Settlement | null>("settlement/" + id);
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
    const response = await axios.post<Settlement>("settlement", request);
    return response.data;
  },
};

export default SettlementApi;
