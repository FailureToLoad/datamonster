import {gql} from '@/__generated__';
import {Survivor, SurvivorGender, SurvivorStatus} from '@types';

export const CreateSurvivor = gql(/* GraphQL */ `
  mutation CreateSurvivor($input: CreateSurvivorInput!) {
    createSurvivor(input: $input) {
      id
      name
    }
  }
`);

export const UpdateSurvivor = gql(/* GraphQL */ `
  mutation UpdateSurvivor($id: ID!, $input: UpdateSurvivorInput!) {
    updateSurvivor(id: $id, input: $input) {
      id
      name
    }
  }
`);

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
  status?: SurvivorStatus;
  statusChangedYear?: number;
};

export type UpdateSurvivorRequest = {
  name?: string;
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
  status?: SurvivorStatus;
  statusChangeYear?: number;
};

export const DefaultSurvivor: Survivor = {
  id: '',
  settlementID: '',
  name: 'Meat',
  born: 0,
  gender: SurvivorGender.M,
  huntxp: 0,
  survival: 1,
  systemicpressure: 0,
  movement: 5,
  accuracy: 0,
  strength: 0,
  evasion: 0,
  luck: 0,
  speed: 0,
  lumi: 0,
  insanity: 0,
  torment: 0,
  courage: 0,
  understanding: 0,
  status: SurvivorStatus.Alive,
  statusChangeYear: 0,
};

export const GET_SURVIVORS = gql(/* GraphQL */ `
  query GetSurvivors($settlementId: ID!) {
    survivors(filter: {settlementID: $settlementId}) {
      id
      accuracy
      born
      courage
      evasion
      gender
      huntxp
      insanity
      luck
      lumi
      movement
      name
      speed
      strength
      survival
      systemicpressure
      torment
      understanding
    }
  }
`);
