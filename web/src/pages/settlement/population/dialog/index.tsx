import {useAppForm} from '~/lib/form';
import {useRef, useLayoutEffect, useState} from 'react';
import {type Survivor, SurvivorGender, SurvivorStatus} from '~/lib/survivor';
import {type} from 'arktype';
import {PostJSON, PatchJSON} from '~/lib/request';
import {GenderFemaleIcon, GenderMaleIcon, DnaIcon, CircleHalfIcon} from '@phosphor-icons/react';
import {StatsTab} from './stats.tsx';
import {TraitsTab} from './traits.tsx';
import {BoxTrack} from '~/components/BoxTrack';
import {SurvivorFormSchema, type SurvivorFormFields, nameValidator, statusValidator} from "./form.ts"
import styles from "./index.module.css";

type SurvivorDialogProps = {
  data: Survivor;
  settlementId: string;
  onClose: () => void;
  onSuccess: () => void;
};

export type SurvivorUpdateRequest = {
  statUpdates?: Record<string, number>;
  statusUpdate?: SurvivorStatus;
  disorders?: string[];
};

export default function SurvivorDialog({data, settlementId, onClose, onSuccess}: SurvivorDialogProps) {
  const dialogRef = useRef<HTMLDialogElement>(null);
  const [cannotSpendSurvival, setCannotSpendSurvival] = useState(false);
  const isCreateMode = !data.id || data.id.trim().length === 0;

  const form = useAppForm({
    defaultValues: {
      name: data.name,
      gender: data.gender,
      status: data.status,
      survival: data.survival,
      systemicPressure: data.systemicPressure,
      movement: data.movement,
      accuracy: data.accuracy,
      strength: data.strength,
      evasion: data.evasion,
      luck: data.luck,
      speed: data.speed,
      lumi: data.lumi,
      insanity: data.insanity,
      torment: data.torment,
      birth: data.birth,
      huntxp: data.huntxp,
      courage: data.courage,
      understanding: data.understanding,
      disorder1: data.disorders?.[0] ?? null,
      disorder2: data.disorders?.[1] ?? null,
      disorder3: data.disorders?.[2] ?? null,
      fightingArt: null,
      secretFightingArt: null,
    } as SurvivorFormFields,
    onSubmit: async ({value}) => {
      const parsed = SurvivorFormSchema(value);
      if (parsed instanceof type.errors) {
        return;
      }

      const disorders = [parsed.disorder1, parsed.disorder2, parsed.disorder3]
        .filter((d): d is string => d !== null);

      if (isCreateMode) {
        const response = await PostJSON(`/api/settlements/${settlementId}/survivors`, {
          name: parsed.name,
          gender: parsed.gender,
          survival: parsed.survival,
          systemicpressure: parsed.systemicPressure,
          movement: parsed.movement,
          accuracy: parsed.accuracy,
          strength: parsed.strength,
          evasion: parsed.evasion,
          luck: parsed.luck,
          speed: parsed.speed,
          lumi: parsed.lumi,
          insanity: parsed.insanity,
          torment: parsed.torment,
          birth: parsed.birth,
          huntxp: parsed.huntxp,
          courage: parsed.courage,
          understanding: parsed.understanding,
          settlementId: settlementId,
          disorders: disorders.length > 0 ? disorders : undefined,
        });

        if (response.ok) {
          dialogRef.current?.close();
          form.reset();
          onSuccess();
        }
      } else {
        const statFields: Array<[keyof SurvivorFormFields, keyof Survivor, string]> = [
          ['huntxp', 'huntxp', 'huntxp'],
          ['survival', 'survival', 'survival'],
          ['movement', 'movement', 'movement'],
          ['accuracy', 'accuracy', 'accuracy'],
          ['strength', 'strength', 'strength'],
          ['evasion', 'evasion', 'evasion'],
          ['luck', 'luck', 'luck'],
          ['speed', 'speed', 'speed'],
          ['lumi', 'lumi', 'lumi'],
          ['insanity', 'insanity', 'insanity'],
          ['torment', 'torment', 'torment'],
          ['systemicPressure', 'systemicPressure', 'systemicPressure'],
          ['courage', 'courage', 'courage'],
          ['understanding', 'understanding', 'understanding'],
        ];

        const statUpdates = statFields.reduce<Record<string, number>>((acc, [formKey, survivorKey, apiKey]) => {
          if (parsed[formKey] !== data[survivorKey]) {
            acc[apiKey] = parsed[formKey] as number;
          }
          return acc;
        }, {});

        const statusUpdate = parsed.status !== data.status ? parsed.status : undefined;

        if (Object.keys(statUpdates).length === 0 && !statusUpdate && disorders.length === 0) {
          dialogRef.current?.close();
          onClose();
          return;
        }

        const payload: SurvivorUpdateRequest = {};
        if (Object.keys(statUpdates).length > 0) {
          payload.statUpdates = statUpdates;
        }
        if (statusUpdate) {
          payload.statusUpdate = statusUpdate;
        }
        if (disorders.length > 0) {
          payload.disorders = disorders;
        }

        const response = await PatchJSON(`/api/settlements/${settlementId}/survivors/${data.id}`, payload);

        if (response.ok) {
          dialogRef.current?.close();
          onSuccess();
        }
      }
    },
  });

  useLayoutEffect(() => {
    dialogRef.current?.showModal();
  }, []);

  useLayoutEffect(() => {
    const dialog = dialogRef.current;
    if (!dialog) return;

    const handleDialogClose = () => {
      onClose();
    };

    dialog.addEventListener('close', handleDialogClose);
    return () => dialog.removeEventListener('close', handleDialogClose);
  }, [onClose]);

  const handleClose = () => {
    dialogRef.current?.close();
  };

  return (
    <dialog ref={dialogRef} className={styles.dialog}>
      <div className={styles.dialogBox}>
        <form
          className={styles.form}
          onSubmit={(e) => {
            e.preventDefault();
            e.stopPropagation();
            form.handleSubmit();
          }}
        >
          <div id="base-details" className={styles.baseDetails}>
            <div className={styles.header}>
              <div className={styles.nameSection}>
                {isCreateMode ? (
                  <form.Field
                    name="name"
                    validators={{
                      onChange: ({value}) => {
                        const result = nameValidator(value);
                        return result instanceof type.errors
                          ? result.summary
                          : undefined;
                      },
                    }}
                  >
                    {(field) => (
                      <input
                        type="text"
                        className={styles.nameInput}
                        value={field.state.value}
                        onChange={(e) => field.handleChange(e.target.value)}
                        placeholder="Survivor name"
                      />
                    )}
                  </form.Field>
                ) : (
                  <p className={styles.name}>{data.name}</p>
                )}
                <form.Field name="gender">
                  {(field) => (
                    isCreateMode ? (
                      <button
                        type="button"
                        className={field.state.value === SurvivorGender.M ? styles.genderBadgeMale : styles.genderBadgeFemale}
                        onClick={() => field.handleChange(
                          field.state.value === SurvivorGender.M ? SurvivorGender.F : SurvivorGender.M
                        )}
                      >
                        {field.state.value === SurvivorGender.M ? <GenderMaleIcon weight="bold" /> : <GenderFemaleIcon weight="bold" />}
                      </button>
                    ) : (
                      <span className={styles.genderBadge}>
                        {data.gender === SurvivorGender.M ? <GenderMaleIcon weight="bold" /> : <GenderFemaleIcon weight="bold" />}
                      </span>
                    )
                  )}
                </form.Field>
                {!isCreateMode && (
                  <span className={styles.birthBadge}>Born year {data.birth}</span>
                )}
              </div>
              {!isCreateMode && (
                <form.Field
                  name="status"
                  validators={{
                    onChange: ({value}) => {
                      const result = statusValidator(value);
                      return result instanceof type.errors
                        ? result.summary
                        : undefined;
                    },
                  }}
                >
                  {(field) => (
                    <select
                      id="status-input"
                      className={styles.statusSelect}
                      value={field.state.value}
                      onChange={(e) =>
                        field.handleChange(e.target.value as SurvivorStatus)
                      }
                    >
                      {Object.values(SurvivorStatus).map((status) => (
                        <option key={status} value={status}>
                          {status}
                        </option>
                      ))}
                    </select>
                  )}
                </form.Field>
              )}
            </div>
            <form.Field name="huntxp">
              {(field) => (
                <BoxTrack
                  value={field.state.value}
                  onChange={(val) => field.handleChange(val)}
                  label="Hunt XP"
                  totalBoxes={16}
                  accentedBoxes={[2, 6, 10, 15, 16]}
                  labelPosition="left"
                />
              )}
            </form.Field>
            <div className={styles.survivalRow}>
              <button
                type="button"
                className={cannotSpendSurvival ? styles.survivalBadgeInactive : styles.survivalBadgeActive}
                onClick={() => setCannotSpendSurvival(!cannotSpendSurvival)}
              >
                {cannotSpendSurvival ? 'Cannot spend survival' : 'Can spend survival'}
              </button>
              <form.Subscribe selector={(state) => state.values.insanity}>
                {(insanity) => insanity >= 3 && (
                  <span className={styles.insaneBadge}>Insane</span>
                )}
              </form.Subscribe>
            </div>
          </div>

          <div className={styles.tabs}>
            <label className={styles.tab}>
              <input type="radio" name="survivor_tabs" defaultChecked />
              <DnaIcon size={20} weight="bold" />
            </label>
            <div className={styles.tabContent}>
              <StatsTab form={form} />
            </div>
            <label className={styles.tab}>
              <input type="radio" name="survivor_tabs" />
              <CircleHalfIcon size={20} weight="fill" />
            </label>
            <div className={styles.tabContent}>
              <TraitsTab form={form} />
            </div>
          </div>

          <div className={styles.actions}>
            <button
              type="button"
              className={styles.cancelBtn}
              onClick={handleClose}
            >
              Cancel
            </button>
            <form.Subscribe selector={(state) => state.canSubmit}>
              {(canSubmit) => (
                <button
                  type="submit"
                  className={styles.saveBtn}
                  disabled={!canSubmit}
                >
                  Save
                </button>
              )}
            </form.Subscribe>
          </div>
        </form>
      </div>
      <form method="dialog" className={styles.backdrop}>
        <button onClick={handleClose}>close</button>
      </form>
    </dialog>
  );
}
