import {useForm} from '@tanstack/react-form';
import {useRef, useLayoutEffect, useState} from 'react';
import {type Survivor, SurvivorGender, SurvivorStatus} from '~/lib/survivor';
import {type} from 'arktype';
import {PostJSON, PatchJSON} from '~/lib/request';
import {BoxTrack} from '~/components/BoxTrack'
import {GenderFemaleIcon, GenderMaleIcon} from '@phosphor-icons/react';
import styles from './survivorDialog.module.css';

const survivorNameValidator = type('1 <= string <= 50');
const isInteger = type('number.integer');
const isPositive = type('number.integer >= 0');
const genderValidator = type.enumerated(SurvivorGender.M, SurvivorGender.F);
const statusValidator = type.enumerated(
  SurvivorStatus.Alive,
  SurvivorStatus.CannotDepart,
  SurvivorStatus.CeasedToExist,
  SurvivorStatus.Dead,
  SurvivorStatus.Retired,
);

const SurvivorSchema = type({
  name: survivorNameValidator,
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
});

type SurvivorFields = typeof SurvivorSchema.infer;

type SurvivorDialogProps = {
  data: Survivor;
  settlementId: string;
  onClose: () => void;
  onSuccess: () => void;
};

export type SurvivorUpdateRequest = {
  statUpdates?: Record<string, number>;
  statusUpdate?: SurvivorStatus;
};

export default function SurvivorDialog({data, settlementId, onClose, onSuccess}: SurvivorDialogProps) {
  const dialogRef = useRef<HTMLDialogElement>(null);
  const [cannotSpendSurvival, setCannotSpendSurvival] = useState(false);
  const isCreateMode = data.id.trim().length === 0;

  const form = useForm({
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
    } as SurvivorFields,
    onSubmit: async ({value}) => {
      const parsed = SurvivorSchema(value);
      if (parsed instanceof type.errors) {
        return;
      }

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
        });

        if (response.ok) {
          dialogRef.current?.close();
          form.reset();
          onSuccess();
        }
      } else {
        const survivor = data as Survivor;

        const statFields: Array<[keyof SurvivorFields, keyof Survivor, string]> = [
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
          if (parsed[formKey] !== survivor[survivorKey]) {
            acc[apiKey] = parsed[formKey] as number;
          }
          return acc;
        }, {});

        const statusUpdate = parsed.status !== survivor.status ? parsed.status : undefined;

        if (Object.keys(statUpdates).length === 0 && !statusUpdate) {
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

        const response = await PatchJSON(`/api/settlements/${survivor.settlementId}/survivors/${survivor.id}`, payload);

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
          onSubmit={(e) => {
            e.preventDefault();
            e.stopPropagation();
            form.handleSubmit();
          }}
        >
          <section className={styles.formGrid}>
            <div id="base-details" className={styles.baseDetails}>
              <div className={styles.nameRow}>
                <div className={styles.nameInfo}>
                  {isCreateMode ? (
                    <form.Field
                      name="name"
                      validators={{
                        onChange: ({value}) => {
                          const result = survivorNameValidator(value);
                          return result instanceof type.errors
                            ? result.summary
                            : undefined;
                        },
                      }}
                    >
                      {(field) => (
                        <input
                          type="text"
                          id="name-input"
                          className={styles.nameInput}
                          value={field.state.value}
                          onChange={(e) => field.handleChange(e.target.value)}
                          placeholder="Survivor name"
                        />
                      )}
                    </form.Field>
                  ) : (
                    <p className={styles.survivorName}>
                      {data.name}
                    </p>
                  )}
                  <form.Field name="gender">
                    {(field) => (
                      isCreateMode ? (
                        <button
                          type="button"
                          className={field.state.value === SurvivorGender.M ? styles.badgeMale : styles.badgeFemale}
                          onClick={() => field.handleChange(
                            field.state.value === SurvivorGender.M ? SurvivorGender.F : SurvivorGender.M
                          )}
                        >
                          {field.state.value === SurvivorGender.M ? <GenderMaleIcon weight="bold" /> : <GenderFemaleIcon weight="bold" />}
                        </button>
                      ) : (
                        <span className={styles.badge}>
                          {data.gender === SurvivorGender.M ? <GenderMaleIcon weight="bold" /> : <GenderFemaleIcon weight="bold" />}
                        </span>
                      )
                    )}
                  </form.Field>
                  {!isCreateMode && (
                    <span className={styles.badgeSemibold}>Born year {data.birth}</span>
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
                        className={styles.select}
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
                  <HuntXPTrack
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                  />
                )}
              </form.Field>
              <div className={styles.badgeRow}>
                <button
                  type="button"
                  className={cannotSpendSurvival ? styles.badgeError : styles.badgeSuccess}
                  onClick={() => setCannotSpendSurvival(!cannotSpendSurvival)}
                >
                  {cannotSpendSurvival ? 'Cannot spend survival' : 'Can spend survival'}
                </button>
                <form.Subscribe selector={(state) => state.values.insanity}>
                  {(insanity) => insanity >= 3 && (
                    <span className={styles.badgeSemibold}>Insane</span>
                  )}
                </form.Subscribe>
              </div>
            </div>

            <div className={styles.statsSection}>
              <form.Field name="survival">
                {(field) => (
                  <StatBox
                    id="survival"
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                    label="Survival"
                  />
                )}
              </form.Field>
              <form.Field name="insanity">
                {(field) => (
                  <StatBox
                    id="insanity"
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                    label="Insanity"
                  />
                )}
              </form.Field>
              <form.Field name="systemicPressure">
                {(field) => (
                  <StatBox
                    id="systemicPressure"
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                    label="Systemic Pressure"
                  />
                )}
              </form.Field>
              <form.Field name="torment">
                {(field) => (
                  <StatBox
                    id="torment"
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                    label="Torment"
                  />
                )}
              </form.Field>
              <form.Field name="lumi">
                {(field) => (
                  <StatBox
                    id="lumi"
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                    label="Lumi"
                  />
                )}
              </form.Field>
            </div>

            <div id="base-stats" className={styles.statsSection}>
              <form.Field name="movement">
                {(field) => (
                  <StatBox
                    id="movement"
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                    label="Movement"
                  />
                )}
              </form.Field>
              <form.Field name="accuracy">
                {(field) => (
                  <StatBox
                    id="accuracy"
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                    label="Accuracy"
                  />
                )}
              </form.Field>
              <form.Field name="strength">
                {(field) => (
                  <StatBox
                    id="strength"
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                    label="Strength"
                  />
                )}
              </form.Field>
              <form.Field name="evasion">
                {(field) => (
                  <StatBox
                    id="evasion"
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                    label="Evasion"
                  />
                )}
              </form.Field>
              <form.Field name="luck">
                {(field) => (
                  <StatBox
                    id="luck"
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                    label="Luck"
                  />
                )}
              </form.Field>
              <form.Field name="speed">
                {(field) => (
                  <StatBox
                    id="speed"
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                    label="Speed"
                  />
                )}
              </form.Field>
            </div>

            <div className={styles.tracksSection}>
              <form.Field name="courage">
                {(field) => (
                  <CourageTrack
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                  />
                )}
              </form.Field>

              <form.Field name="understanding">
                {(field) => (
                  <UnderstandingTrack
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                  />
                )}
              </form.Field>
            </div>
          </section>

          <div className={styles.modalFooter}>
            <button
              type="button"
              className={styles.btn}
              onClick={handleClose}
            >
              Cancel
            </button>
            <form.Subscribe selector={(state) => state.canSubmit}>
              {(canSubmit) => (
                <button
                  type="submit"
                  className={styles.btnPrimary}
                  disabled={!canSubmit}
                >
                  {isCreateMode ? 'Create' : 'Save'}
                </button>
              )}
            </form.Subscribe>
          </div>
        </form>
      </div>
      <form method="dialog" className="modal-backdrop">
        <button onClick={handleClose}>close</button>
      </form>
    </dialog>
  );
}

function StatBox({
  id,
  value,
  onChange,
  label,
}: {
  id: string;
  value: number;
  onChange: (val: number) => void;
  label: string;
}) {
  return (
    <div className={`${styles.statBox} ${styles.statBoxFlex1}`}>
      <div className={styles.statBoxInput}>
        <input
          id={`${id}-input`}
          type="number"
          className={styles.statBoxInputField}
          value={value}
          onChange={(e) => onChange(Number(e.target.value))}
        />
      </div>
      <p className={styles.statBoxLabel}>
        {label}
      </p>
    </div>
  );
}


function HuntXPTrack({
  value,
  onChange,
}: {
  value: number;
  onChange: (val: number) => void;
}) {
  return (
    <BoxTrack
      value={value}
      onChange={onChange}
      label="Hunt XP"
      totalBoxes={16}
      accentedBoxes={[2, 6, 10, 15, 16]}
      labelPosition="left"
    />
  );
}

function CourageTrack({
  value,
  onChange,
}: {
  value: number;
  onChange: (val: number) => void;
}) {
  return (
    <BoxTrack
      value={value}
      onChange={onChange}
      label="Courage"
      totalBoxes={9}
      accentedBoxes={[3, 9]}
      labelPosition="top"
      className={styles.flex1}
    />
  );
}

function UnderstandingTrack({
  value,
  onChange,
}: {
  value: number;
  onChange: (val: number) => void;
}) {
  return (
    <BoxTrack
      value={value}
      onChange={onChange}
      label="Understanding"
      totalBoxes={9}
      accentedBoxes={[3, 9]}
      labelPosition="top"
      className={styles.flex1}
    />
  );
}
