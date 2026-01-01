import {useForm} from '@tanstack/react-form';
import {useRef} from 'react';
import {PlusIcon} from '@phosphor-icons/react';
import {SurvivorGender} from '~/lib/survivor';
import {type} from 'arktype';
import { PostJSON } from '~/lib/request';
import styles from './survivorDialog.module.css';

const survivorNameValidator = type('1 <= string <= 50');
const isInteger = type('number.integer');
const isPositive = type('number.integer >= 0');
const genderValidator = type.enumerated(SurvivorGender.M, SurvivorGender.F);

const AddSurvivorSchema = type({
  name: survivorNameValidator,
  gender: genderValidator,
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
});

type AddSurvivorFields = typeof AddSurvivorSchema.infer;

type SurvivorDialogProps = {
  settlementId: string;
  onSuccess?: () => void;
};

export default function NewSurvivorDialog({settlementId, onSuccess}: SurvivorDialogProps) {
  const dialogRef = useRef<HTMLDialogElement>(null);

  const form = useForm({
    defaultValues: {
      name: 'Meat',
      gender: SurvivorGender.M,
      survival: 1,
      systemicPressure: 0,
      movement: 5,
      accuracy: 0,
      strength: 0,
      evasion: 0,
      luck: 0,
      speed: 0,
      lumi: 0,
      insanity: 0,
      torment: 0,
    } as AddSurvivorFields,
    onSubmit: async ({value}) => {
      const parsed = AddSurvivorSchema(value);
      if (parsed instanceof type.errors) {
        return;
      }

      const response = await PostJSON(`/api/settlements/${settlementId}/survivors`, 
        {
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
          birth: 1,
          huntxp: 0,
          courage: 0,
          understanding: 0,
          settlementId: settlementId,
        })
      if (response.ok) {
        dialogRef.current?.close();
        form.reset();
        onSuccess?.();
      }
    },
  });

  return (
    <>
      <button
        className={styles.btnGhost}
        aria-label="Create Survivor"
        title="Create Survivor"
        onClick={() => dialogRef.current?.showModal()}
      >
        <PlusIcon size={18} weight="bold" />
      </button>
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
              {/* Section 1: Name + Gender row */}
              <div className={styles.nameRow}>
                <div className={styles.nameInputWrapper}>
                  <p className={styles.sectionLabel}>
                    Name
                  </p>
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
                      />
                    )}
                  </form.Field>
                </div>
                <form.Field name="gender">
                  {(field) => (
                    <div className={styles.genderGroup}>
                      <label className={styles.radioLabel}>
                        <input
                          type="radio"
                          name="gender"
                          className={styles.radio}
                          value="M"
                          checked={field.state.value === 'M'}
                          onChange={() => field.handleChange(SurvivorGender.M)}
                        />
                        <span>M</span>
                      </label>
                      <label className={styles.radioLabel}>
                        <input
                          type="radio"
                          name="gender"
                          className={styles.radio}
                          value="F"
                          checked={field.state.value === 'F'}
                          onChange={() => field.handleChange(SurvivorGender.F)}
                        />
                        <span>F</span>
                      </label>
                    </div>
                  )}
                </form.Field>
              </div>

              {/* Section 2: Survival box */}
              <div className={styles.survivalSection}>
                <form.Field
                  name="survival"
                  validators={{
                    onChange: ({value}) => {
                      const result = isPositive(value);
                      return result instanceof type.errors
                        ? result.summary
                        : undefined;
                    },
                  }}
                >
                  {(field) => (
                    <div className={styles.survivalBox}>
                      <input
                        id="survival-input"
                        type="number"
                        className={styles.survivalInput}
                        value={field.state.value}
                        onChange={(e) =>
                          field.handleChange(Number(e.target.value))
                        }
                      />
                    </div>
                  )}
                </form.Field>

                <div className={styles.survivalInfo}>
                  <p className={styles.sectionLabel}>
                    Survival
                  </p>
                  <div className={styles.checkboxRow}>
                    <div className={styles.checkboxItem}>
                      <input
                        type="checkbox"
                        id="cannot-spend"
                        className={styles.checkbox}
                      />
                      <label
                        htmlFor="cannot-spend"
                        className={styles.checkboxLabel}
                      >
                        Cannot spend survival
                      </label>
                    </div>
                  </div>
                </div>
                <div className={styles.skillsColumn}>
                  <CheckboxItem id="dodge" label="Dodge" />
                  <CheckboxItem id="encourage" label="Encourage" />
                  <CheckboxItem id="surge" label="Surge" />
                  <CheckboxItem id="dash" label="Dash" />
                  <CheckboxItem id="fistpump" label="Fist Pump" />
                </div>
                <div className={styles.systemicPressureWrapper}>
                  <form.Field name="systemicPressure">
                    {(field) => (
                      <StatBox
                        id="systemicPressure"
                        value={field.state.value}
                        onChange={(val) => field.handleChange(val)}
                        label="Systemic Pressure"
                        className={styles.systemicPressureStat}
                      />
                    )}
                  </form.Field>
                </div>
              </div>

              {/* Section 3: Stats row */}
              <div className={styles.statsRow}>
                <form.Field name="movement">
                  {(field) => (
                    <StatBox
                      id="movement"
                      value={field.state.value}
                      onChange={(val) => field.handleChange(val)}
                      label="Movement"
                      className={styles.statFirst}
                    />
                  )}
                </form.Field>
                <div className={styles.divider} />
                <form.Field name="accuracy">
                  {(field) => (
                    <StatBox
                      id="accuracy"
                      value={field.state.value}
                      onChange={(val) => field.handleChange(val)}
                      label="Accuracy"
                      className={styles.statMiddle}
                    />
                  )}
                </form.Field>
                <div className={styles.divider} />
                <form.Field name="strength">
                  {(field) => (
                    <StatBox
                      id="strength"
                      value={field.state.value}
                      onChange={(val) => field.handleChange(val)}
                      label="Strength"
                      className={styles.statMiddle}
                    />
                  )}
                </form.Field>
                <div className={styles.divider} />
                <form.Field name="evasion">
                  {(field) => (
                    <StatBox
                      id="evasion"
                      value={field.state.value}
                      onChange={(val) => field.handleChange(val)}
                      label="Evasion"
                      className={styles.statMiddle}
                    />
                  )}
                </form.Field>
                <div className={styles.divider} />
                <form.Field name="luck">
                  {(field) => (
                    <StatBox
                      id="luck"
                      value={field.state.value}
                      onChange={(val) => field.handleChange(val)}
                      label="Luck"
                      className={styles.statMiddle}
                    />
                  )}
                </form.Field>
                <div className={styles.divider} />
                <form.Field name="speed">
                  {(field) => (
                    <StatBox
                      id="speed"
                      value={field.state.value}
                      onChange={(val) => field.handleChange(val)}
                      label="Speed"
                      className={styles.statMiddle}
                    />
                  )}
                </form.Field>
                <div className={styles.divider} />
                <form.Field name="lumi">
                  {(field) => (
                    <StatBox
                      id="lumi"
                      value={field.state.value}
                      onChange={(val) => field.handleChange(val)}
                      label="Lumi"
                      className={styles.statLast}
                    />
                  )}
                </form.Field>
              </div>

              {/* Section 4: Brain row */}
              <div className={styles.brainSection}>
                <form.Field name="insanity">
                  {(field) => (
                    <StatBox
                      id="insanity"
                      value={field.state.value}
                      onChange={(val) => field.handleChange(val)}
                      label="Insanity"
                      className={styles.statBoxOrderFirst}
                    />
                  )}
                </form.Field>
                <div className={styles.divider} />
                <div className={styles.brainInfo}>
                  <div className={styles.brainHeader}>
                    <p className={styles.sectionLabel}>
                      Brain
                    </p>
                    <input
                      type="checkbox"
                      id="brainbox"
                      className={styles.brainCheckbox}
                    />
                  </div>

                  <div className={styles.checkboxRow}>
                    <p className={styles.brainDescription}>
                      If your insanity is 3+, you are <b>insane</b>
                    </p>
                  </div>
                </div>
                <div className={styles.divider} />
                <form.Field name="torment">
                  {(field) => (
                    <StatBox
                      id="torment"
                      value={field.state.value}
                      onChange={(val) => field.handleChange(val)}
                      label="Torment"
                      className={styles.statBoxOrderLast}
                    />
                  )}
                </form.Field>
              </div>
            </section>

            {/* Footer */}
            <div className={styles.modalFooter}>
              <button
                type="button"
                className={styles.btn}
                onClick={() => dialogRef.current?.close()}
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
                    Create
                  </button>
                )}
              </form.Subscribe>
            </div>
          </form>
        </div>
        <form method="dialog" className="modal-backdrop">
          <button>close</button>
        </form>
      </dialog>
    </>
  );
}

function CheckboxItem({id, label}: {id: string; label: string}) {
  return (
    <div className={styles.checkboxItem}>
      <input type="checkbox" id={id} className={styles.checkboxSmall} />
      <label
        htmlFor={id}
        className={styles.checkboxLabel}
      >
        {label}
      </label>
    </div>
  );
}

function StatBox({
  id,
  value,
  onChange,
  label,
  className = '',
}: {
  id: string;
  value: number;
  onChange: (val: number) => void;
  label: string;
  className?: string;
}) {
  return (
    <div
      className={`${styles.statBox} ${className}`}
    >
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
