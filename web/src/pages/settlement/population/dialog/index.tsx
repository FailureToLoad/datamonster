import { createForm } from '@tanstack/solid-form';
import { onMount, onCleanup, Show, For } from 'solid-js';
import { type Survivor, SurvivorGender, SurvivorStatus } from '~/lib/survivor';
import { type } from 'arktype';
import { PostJSON, PatchJSON } from '~/lib/request';
import Dna from 'lucide-solid/icons/dna';
import Contrast from 'lucide-solid/icons/contrast';
import { StatsTab } from './stats';
import { TraitsTab } from './traits';
import { BoxTrack } from '~/components/BoxTrack';
import { SurvivorFormSchema, type SurvivorFormFields, nameValidator, statusValidator } from './form';
import styles from './index.module.css';

interface SurvivorDialogProps {
  data: Survivor;
  settlementId: string;
  onClose: () => void;
  onSuccess: () => void;
}

export type SurvivorUpdateRequest = {
  statUpdates?: Record<string, number>;
  statusUpdate?: SurvivorStatus;
  disorders: string[];
  fightingArt: string | null;
  secretFightingArt: string | null;
};

export default function SurvivorDialog(props: SurvivorDialogProps) {
  let dialogRef: HTMLDialogElement | undefined;
  const isCreateMode = () => !props.data.id || props.data.id.trim().length === 0;

  const form = createForm(() => ({
    defaultValues: {
      name: props.data.name,
      gender: props.data.gender,
      status: props.data.status,
      survival: props.data.survival,
      systemicPressure: props.data.systemicPressure,
      movement: props.data.movement,
      accuracy: props.data.accuracy,
      strength: props.data.strength,
      evasion: props.data.evasion,
      luck: props.data.luck,
      speed: props.data.speed,
      lumi: props.data.lumi,
      insanity: props.data.insanity,
      torment: props.data.torment,
      birth: props.data.birth,
      huntxp: props.data.huntxp,
      courage: props.data.courage,
      understanding: props.data.understanding,
      disorder1: props.data.disorders?.[0] ?? null,
      disorder2: props.data.disorders?.[1] ?? null,
      disorder3: props.data.disorders?.[2] ?? null,
      fightingArt: props.data.fightingArt ?? null,
      secretFightingArt: props.data.secretFightingArt ?? null,
    } as SurvivorFormFields,
    onSubmit: async ({ value }) => {
      const parsed = SurvivorFormSchema(value);
      if (parsed instanceof type.errors) {
        return;
      }

      const disorders = [parsed.disorder1, parsed.disorder2, parsed.disorder3].filter(
        (d): d is string => d !== null
      );

      if (isCreateMode()) {
        const response = await PostJSON(`/api/settlements/${props.settlementId}/survivors`, {
          name: parsed.name,
          gender: parsed.gender,
          survival: parsed.survival,
          systemicPressure: parsed.systemicPressure,
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
          settlementId: props.settlementId,
          disorders: disorders.length > 0 ? disorders : undefined,
          fightingArt: parsed.fightingArt ?? undefined,
          secretFightingArt: parsed.secretFightingArt ?? undefined,
        });

        if (response.ok) {
          dialogRef?.close();
          form.reset();
          props.onSuccess();
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
          if (parsed[formKey] !== props.data[survivorKey]) {
            acc[apiKey] = parsed[formKey] as number;
          }
          return acc;
        }, {});

        const statusUpdate = parsed.status !== props.data.status ? parsed.status : undefined;

        const payload: SurvivorUpdateRequest = {
          disorders,
          fightingArt: parsed.fightingArt,
          secretFightingArt: parsed.secretFightingArt,
        };
        if (Object.keys(statUpdates).length > 0) {
          payload.statUpdates = statUpdates;
        }
        if (statusUpdate) {
          payload.statusUpdate = statusUpdate;
        }

        const response = await PatchJSON(`/api/settlements/${props.settlementId}/survivors/${props.data.id}`, payload);

        if (response.ok) {
          dialogRef?.close();
          props.onSuccess();
        }
      }
    },
  }));

  onMount(() => {
    dialogRef?.showModal();
  });

  onMount(() => {
    const dialog = dialogRef;
    if (!dialog) return;

    const handleDialogClose = () => {
      props.onClose();
    };

    dialog.addEventListener('close', handleDialogClose);
    onCleanup(() => dialog.removeEventListener('close', handleDialogClose));
  });

  const handleClose = () => {
    dialogRef?.close();
  };

  return (
    <dialog ref={dialogRef} class={styles.dialog}>
      <div class={styles.dialogBox}>
        <form
          class={styles.form}
          onSubmit={(e) => {
            e.preventDefault();
            e.stopPropagation();
            form.handleSubmit();
          }}
        >
          <div id="base-details" class={styles.baseDetails}>
            <div class={styles.header}>
              <div class={styles.nameSection}>
                <Show
                  when={isCreateMode()}
                  fallback={<p class={styles.name}>{props.data.name}</p>}
                >
                  <form.Field
                    name="name"
                    validators={{
                      onChange: ({ value }) => {
                        const result = nameValidator(value);
                        return result instanceof type.errors ? result.summary : undefined;
                      },
                    }}
                  >
                    {(field) => (
                      <input
                        type="text"
                        class={styles.nameInput}
                        value={field().state.value}
                        onInput={(e) => field().handleChange(e.currentTarget.value)}
                        placeholder="Survivor name"
                      />
                    )}
                  </form.Field>
                </Show>
                <form.Field name="gender">
                  {(field) => (
                    <Show
                      when={isCreateMode()}
                      fallback={
                        <span class={styles.genderBadge}>
                          {props.data.gender === SurvivorGender.M ? 'M' : 'F'}
                        </span>
                      }
                    >
                      <button
                        type="button"
                        class={field().state.value === SurvivorGender.M ? styles.genderBadgeMale : styles.genderBadgeFemale}
                        onClick={() =>
                          field().handleChange(field().state.value === SurvivorGender.M ? SurvivorGender.F : SurvivorGender.M)
                        }
                      >
                        {field().state.value === SurvivorGender.M ? 'M' : 'F'}
                      </button>
                    </Show>
                  )}
                </form.Field>
                <Show when={!isCreateMode()}>
                  <span class={styles.birthBadge}>Born year {props.data.birth}</span>
                </Show>
              </div>
              <Show when={!isCreateMode()}>
                <form.Field
                  name="status"
                  validators={{
                    onChange: ({ value }) => {
                      const result = statusValidator(value);
                      return result instanceof type.errors ? result.summary : undefined;
                    },
                  }}
                >
                  {(field) => (
                    <select
                      id="status-input"
                      class={styles.statusSelect}
                      value={field().state.value}
                      onChange={(e) => field().handleChange(e.currentTarget.value as SurvivorStatus)}
                    >
                      <For each={Object.values(SurvivorStatus)}>
                        {(status) => <option value={status}>{status}</option>}
                      </For>
                    </select>
                  )}
                </form.Field>
              </Show>
            </div>
            <form.Field name="huntxp">
              {(field) => (
                <BoxTrack
                  value={field().state.value}
                  onChange={(val) => field().handleChange(val)}
                  label="Hunt XP"
                  totalBoxes={16}
                  accentedBoxes={[2, 6, 10, 15, 16]}
                  labelPosition="left"
                />
              )}
            </form.Field>
            <div class={styles.survivalRow}>
              <form.Subscribe selector={(state) => state.values.insanity}>
                {(insanity) => (
                  <Show when={insanity() >= 3}>
                    <span class={styles.insaneBadge}>Insane</span>
                  </Show>
                )}
              </form.Subscribe>
            </div>
          </div>

          <div class={styles.tabs}>
            <label class={styles.tab}>
              <input type="radio" name="survivor_tabs" checked />
              <Dna size={20} strokeWidth={2.5} />
            </label>
            <div class={styles.tabContent}>
              <StatsTab form={form} />
            </div>
            <label class={styles.tab}>
              <input type="radio" name="survivor_tabs" />
              <Contrast size={20} strokeWidth={2.5} />
            </label>
            <div class={styles.tabContent}>
              <TraitsTab form={form} />
            </div>
          </div>

          <div class={styles.actions}>
            <button type="button" class={styles.cancelBtn} onClick={handleClose}>
              Cancel
            </button>
            <form.Subscribe selector={(state) => state.canSubmit}>
              {(canSubmit) => (
                <button type="submit" class={styles.saveBtn} disabled={!canSubmit()}>
                  Save
                </button>
              )}
            </form.Subscribe>
          </div>
        </form>
      </div>
      <form method="dialog" class={styles.backdrop}>
        <button onClick={handleClose}>close</button>
      </form>
    </dialog>
  );
}
