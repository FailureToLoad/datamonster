import type { UUID } from "node:crypto";

export type Settlement =  SettlementId & {
  survivalLimit: number;
  departingSurvival: number;
  collectiveCognition: number;
  currentYear: number;
};

export type SettlementId = {
  id: UUID;
  name: string;
};