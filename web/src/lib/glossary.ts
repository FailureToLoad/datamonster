export type Disorder = {
  id: string;
  name: string;
  source: string;
  flavorText?: string;
  effect: string;
};

export type FightingArt = {
  id: string;
  name: string;
  secret: boolean;
  source: string;
  text: string[];
  notes?: string;
};
