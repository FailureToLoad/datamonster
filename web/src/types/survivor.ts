
export type Survivor = {
  id: string;
  name: string;
  gender: SurvivorGender;
  born: number;
  huntxp: number;
  movement: number;
  speed: number;
  strength: number;
  accuracy: number;
  evasion: number;
  luck: number;
  systemicpressure: number;
  torment: number;
  courage: number;
  understanding: number;
  survival: number;
  insanity: number;
  lumi: number;
  settlementID: string;
};


export const SurvivorGender = {
  F: 'F',
  M: 'M'
} as const;

export type SurvivorGender = typeof SurvivorGender[keyof typeof SurvivorGender];


export type CreateSurvivorRequest = {
  settlementID: string;
  name: string;
  born?: number;
  gender?: SurvivorGender;
  huntxp?: number;
  survival?: number;
  movement?: number;
  accuracy?: number;
  strength?: number;
  evasion?: number;
  luck?: number;
  speed?: number;
  insanity?: number;
  systemicpressure?: number;
  torment?: number;
  lumi?: number;
  courage?: number;
  understanding?: number;
};