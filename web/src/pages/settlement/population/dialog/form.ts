import { SurvivorGender, SurvivorStatus } from "~/lib/survivor";
import { type } from "arktype";

const isInteger = type("number.integer");
const isPositive = type("number.integer >= 0");
export const nameValidator = type("1 <= string <= 50");
export const genderValidator = type.enumerated(
  SurvivorGender.M,
  SurvivorGender.F
);
export const statusValidator = type.enumerated(
  SurvivorStatus.Alive,
  SurvivorStatus.CannotDepart,
  SurvivorStatus.CeasedToExist,
  SurvivorStatus.Dead,
  SurvivorStatus.Retired
);

const optionalString = type("string | null");

export const SurvivorFormSchema = type({
  name: nameValidator,
  gender: genderValidator,
  status: statusValidator,
  survival: isPositive,
  systemicPressure: isInteger,
  movement: isInteger,
  accuracy: isInteger,
  strength: isInteger,
  evasion: isInteger,
  luck: isInteger,
  speed: isInteger,
  lumi: isPositive,
  insanity: isPositive,
  torment: isInteger,
  birth: isPositive,
  huntxp: isPositive,
  courage: isPositive,
  understanding: isPositive,
  disorder1: optionalString,
  disorder2: optionalString,
  disorder3: optionalString,
  fightingArt: optionalString,
  secretFightingArt: optionalString,
});

export type SurvivorFormFields = typeof SurvivorFormSchema.infer;
