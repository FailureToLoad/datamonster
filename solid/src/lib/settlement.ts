export type Settlement = SettlementId & {
  survivalLimit: number;
  departingSurvival: number;
  collectiveCognition: number;
  currentYear: number;
};

export type SettlementId = {
  id: string;
  name: string;
};
