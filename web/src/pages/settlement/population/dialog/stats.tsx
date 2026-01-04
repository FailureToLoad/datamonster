import {BoxTrack} from '~/components/BoxTrack';
import {withForm} from '~/lib/form';
import type {SurvivorFormFields} from './form';
import styles from './stats.module.css';

export const StatsTab = withForm({
  defaultValues: {} as SurvivorFormFields,
  render: function Render({form}) {
    return (
      <section className={styles.section}>
        <div className={styles.statRow}>
          <form.Field name="survival">
            {(field) => (
              <StatBox
                id="survival-edit"
                value={field.state.value}
                onChange={(val) => field.handleChange(val)}
                label="Survival"
              />
            )}
          </form.Field>
          <form.Field name="insanity">
            {(field) => (
              <StatBox
                id="insanity-edit"
                value={field.state.value}
                onChange={(val) => field.handleChange(val)}
                label="Insanity"
              />
            )}
          </form.Field>
          <form.Field name="systemicPressure">
            {(field) => (
              <StatBox
                id="systemicPressure-edit"
                value={field.state.value}
                onChange={(val) => field.handleChange(val)}
                label="Systemic Pressure"
              />
            )}
          </form.Field>
          <form.Field name="torment">
            {(field) => (
              <StatBox
                id="torment-edit"
                value={field.state.value}
                onChange={(val) => field.handleChange(val)}
                label="Torment"
              />
            )}
          </form.Field>
          <form.Field name="lumi">
            {(field) => (
              <StatBox
                id="lumi-edit"
                value={field.state.value}
                onChange={(val) => field.handleChange(val)}
                label="Lumi"
              />
            )}
          </form.Field>
        </div>

        <div id="base-stats" className={styles.statRow}>
          <form.Field name="movement">
            {(field) => (
              <StatBox
                id="movement-edit"
                value={field.state.value}
                onChange={(val) => field.handleChange(val)}
                label="Movement"
              />
            )}
          </form.Field>
          <form.Field name="accuracy">
            {(field) => (
              <StatBox
                id="accuracy-edit"
                value={field.state.value}
                onChange={(val) => field.handleChange(val)}
                label="Accuracy"
              />
            )}
          </form.Field>
          <form.Field name="strength">
            {(field) => (
              <StatBox
                id="strength-edit"
                value={field.state.value}
                onChange={(val) => field.handleChange(val)}
                label="Strength"
              />
            )}
          </form.Field>
          <form.Field name="evasion">
            {(field) => (
              <StatBox
                id="evasion-edit"
                value={field.state.value}
                onChange={(val) => field.handleChange(val)}
                label="Evasion"
              />
            )}
          </form.Field>
          <form.Field name="luck">
            {(field) => (
              <StatBox
                id="luck-edit"
                value={field.state.value}
                onChange={(val) => field.handleChange(val)}
                label="Luck"
              />
            )}
          </form.Field>
          <form.Field name="speed">
            {(field) => (
              <StatBox
                id="speed-edit"
                value={field.state.value}
                onChange={(val) => field.handleChange(val)}
                label="Speed"
              />
            )}
          </form.Field>
        </div>

        <div className={styles.trackRow}>
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
    );
  },
});

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
    <div className={styles.statBox}>
      <div className={styles.statInputWrapper}>
        <input
          id={`${id}-input`}
          type="number"
          className={styles.statInput}
          value={value}
          onChange={(e) => onChange(Number(e.target.value))}
        />
      </div>
      <p className={styles.statLabel}>
        {label}
      </p>
    </div>
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
    />
  );
}
