import { BoxTrack } from '~/components/BoxTrack';
import type { FormApi } from '@tanstack/solid-form';
import type { SurvivorFormFields } from './form';
import styles from './stats.module.css';

interface StatsTabProps {
  form: FormApi<SurvivorFormFields, undefined>;
}

export function StatsTab(props: StatsTabProps) {
  return (
    <section class={styles.section}>
      <div class={styles.statRow}>
        <props.form.Field name="survival">
          {(field) => (
            <StatBox
              id="survival-edit"
              value={field().state.value}
              onChange={(val) => field().handleChange(val)}
              label="Survival"
            />
          )}
        </props.form.Field>
        <props.form.Field name="insanity">
          {(field) => (
            <StatBox
              id="insanity-edit"
              value={field().state.value}
              onChange={(val) => field().handleChange(val)}
              label="Insanity"
            />
          )}
        </props.form.Field>
        <props.form.Field name="systemicPressure">
          {(field) => (
            <StatBox
              id="systemicPressure-edit"
              value={field().state.value}
              onChange={(val) => field().handleChange(val)}
              label="Systemic Pressure"
            />
          )}
        </props.form.Field>
        <props.form.Field name="torment">
          {(field) => (
            <StatBox
              id="torment-edit"
              value={field().state.value}
              onChange={(val) => field().handleChange(val)}
              label="Torment"
            />
          )}
        </props.form.Field>
        <props.form.Field name="lumi">
          {(field) => (
            <StatBox
              id="lumi-edit"
              value={field().state.value}
              onChange={(val) => field().handleChange(val)}
              label="Lumi"
            />
          )}
        </props.form.Field>
      </div>

      <div id="base-stats" class={styles.statRow}>
        <props.form.Field name="movement">
          {(field) => (
            <StatBox
              id="movement-edit"
              value={field().state.value}
              onChange={(val) => field().handleChange(val)}
              label="Movement"
            />
          )}
        </props.form.Field>
        <props.form.Field name="accuracy">
          {(field) => (
            <StatBox
              id="accuracy-edit"
              value={field().state.value}
              onChange={(val) => field().handleChange(val)}
              label="Accuracy"
            />
          )}
        </props.form.Field>
        <props.form.Field name="strength">
          {(field) => (
            <StatBox
              id="strength-edit"
              value={field().state.value}
              onChange={(val) => field().handleChange(val)}
              label="Strength"
            />
          )}
        </props.form.Field>
        <props.form.Field name="evasion">
          {(field) => (
            <StatBox
              id="evasion-edit"
              value={field().state.value}
              onChange={(val) => field().handleChange(val)}
              label="Evasion"
            />
          )}
        </props.form.Field>
        <props.form.Field name="luck">
          {(field) => (
            <StatBox
              id="luck-edit"
              value={field().state.value}
              onChange={(val) => field().handleChange(val)}
              label="Luck"
            />
          )}
        </props.form.Field>
        <props.form.Field name="speed">
          {(field) => (
            <StatBox
              id="speed-edit"
              value={field().state.value}
              onChange={(val) => field().handleChange(val)}
              label="Speed"
            />
          )}
        </props.form.Field>
      </div>

      <div class={styles.trackRow}>
        <props.form.Field name="courage">
          {(field) => (
            <BoxTrack
              value={field().state.value}
              onChange={(val) => field().handleChange(val)}
              label="Courage"
              totalBoxes={9}
              accentedBoxes={[3, 9]}
              labelPosition="top"
            />
          )}
        </props.form.Field>

        <props.form.Field name="understanding">
          {(field) => (
            <BoxTrack
              value={field().state.value}
              onChange={(val) => field().handleChange(val)}
              label="Understanding"
              totalBoxes={9}
              accentedBoxes={[3, 9]}
              labelPosition="top"
            />
          )}
        </props.form.Field>
      </div>
    </section>
  );
}

function StatBox(props: { id: string; value: number; onChange: (val: number) => void; label: string }) {
  return (
    <div class={styles.statBox}>
      <div class={styles.statInputWrapper}>
        <input
          id={`${props.id}-input`}
          type="number"
          class={styles.statInput}
          value={props.value}
          onInput={(e) => props.onChange(Number(e.currentTarget.value))}
        />
      </div>
      <p class={styles.statLabel}>{props.label}</p>
    </div>
  );
}
