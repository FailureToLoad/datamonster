export type Settlement = {
  id: string;
  name: string;
  limit: number;
  departing: number;
  cc: number;
  year: number;
};

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
