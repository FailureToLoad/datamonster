export type Survivor = {
  id: string;
  name: string;
  gender: SurvivorGender;
  birth: number;
  status: SurvivorStatus;
  huntxp: number;
  movement: number;
  speed: number;
  strength: number;
  accuracy: number;
  evasion: number;
  luck: number;
  systemicPressure: number;
  torment: number;
  courage: number;
  understanding: number;
  survival: number;
  insanity: number;
  lumi: number;
  settlementId: string;
  disorders: string[];
  fightingArt: string | null;
  secretFightingArt: string | null;
};

export const SurvivorGender = {
  F: "F",
  M: "M",
} as const;

export type SurvivorGender =
  (typeof SurvivorGender)[keyof typeof SurvivorGender];

export const SurvivorStatus = {
  Alive: "Alive",
  CeasedToExist: "Ceased to exist",
  CannotDepart: "Cannot depart",
  Dead: "Dead",
  Retired: "Retired",
} as const;

export type SurvivorStatus =
  (typeof SurvivorStatus)[keyof typeof SurvivorStatus];

export function SurvivorTemplate(settlementId: string): Survivor {
  return {
    id: "",
    settlementId: settlementId,
    name: "Anon",
    gender: SurvivorGender.F,
    birth: 1,
    status: SurvivorStatus.Alive,
    huntxp: 0,
    survival: 1,
    movement: 5,
    accuracy: 0,
    strength: 0,
    evasion: 0,
    luck: 0,
    speed: 0,
    systemicPressure: 0,
    torment: 0,
    courage: 0,
    understanding: 0,
    insanity: 0,
    lumi: 0,
    disorders: [],
    fightingArt: null,
    secretFightingArt: null,
  };
}

export type CreateSurvivorRequest = {
  settlementId: string;
  name: string;
  birth?: number;
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
