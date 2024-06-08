import api from "./api";

export type Survivor = {
  id: number;
  settlement: number;
  name: string;
  born: number;
  gender: string;
  status: "alive" | "dead" | "retired";
  huntXp: number;
  survival: number;
  movement: number;
  accuracy: number;
  strength: number;
  evasion: number;
  luck: number;
  speed: number;
  insanity: number;
  systemicPressure: number;
  torment: number;
  lumi: number;
  courage: number;
  understanding: number;
};

type survivorRequests = {
  getSurvivorsForSettlement: (settlementId: string) => Promise<Array<Survivor>>;
  createSurvivor: (
    settlementId: string,
    survivor: Survivor,
  ) => Promise<Survivor>;
};

const SurvivorApi: survivorRequests = {
  getSurvivorsForSettlement: async function (
    settlementId: string,
  ): Promise<Survivor[]> {
    const response = await api.get<Array<Survivor>>(
      `http://dev.local:8080/settlement/${settlementId}/survivor`,
    );
    return response.data;
  },
  createSurvivor: async function (
    settlementId: string,
    survivor: Survivor,
  ): Promise<Survivor> {
    const response = await api.post<Survivor>(
      `http://dev.local:8080/settlement/${settlementId}/survivor`,
      survivor,
    );
    return response.data;
  },
};

export default SurvivorApi;
