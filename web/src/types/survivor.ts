
export type Survivor = {
  accuracy: number;
  born: number;
  courage: number;
  evasion: number;
  gender: SurvivorGender;
  huntxp: number;
  id: string;
  insanity: number;
  luck: number;
  lumi: number;
  movement: number;
  name: string;
  settlementID: string;
  speed: number;
  strength: number;
  survival: number;
  systemicpressure: number;
  torment: number;
  understanding: number;
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