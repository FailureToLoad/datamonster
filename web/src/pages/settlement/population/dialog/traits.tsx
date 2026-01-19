import type { FormApi } from '@tanstack/solid-form';
import type { SurvivorFormFields } from './form';
import styles from './traits.module.css';
import { GlossaryAutocomplete } from './glossaryAutocomplete';
import { useGlossary } from '~/hooks/glossary';
import type { Disorder, FightingArt } from '~/lib/glossary';
import type { JSX } from 'solid-js';

interface TraitsTabProps {
  form: FormApi<SurvivorFormFields, undefined>;
}

function renderDisorder(disorder: Disorder): JSX.Element {
  return (
    <>
      <span class={styles.itemName}>{disorder.name}</span>
      <span class={styles.itemEffect}> - {disorder.effect}</span>
    </>
  );
}

function renderFightingArt(art: FightingArt): JSX.Element {
  return (
    <>
      <span class={styles.itemName}>{art.name}</span>
      <span class={styles.itemEffect}> - {art.text.join(' ')}</span>
    </>
  );
}

export function TraitsTab(props: TraitsTabProps) {
  const glossary = useGlossary();
  const disorders = () => glossary().disorders;
  const fightingArts = () => glossary().fightingArts;

  return (
    <section class={styles.section}>
      <props.form.Subscribe selector={(state) => [state.values.disorder1, state.values.disorder2, state.values.disorder3]}>
        {(disorderValues) => {
          const [d1, d2, d3] = disorderValues();
          const count = [d1, d2, d3].filter(Boolean).length;
          return (
            <details class={styles.card}>
              <summary class="collapse-title">
                <span class={styles.cardTitle}>Disorders</span>
                {count > 0 && <span class={styles.countBadge}>{count}</span>}
              </summary>
              <div class="collapse-content">
                <div class={styles.cardContent}>
                  <props.form.Field name="disorder1">
                    {(field) => (
                      <GlossaryAutocomplete
                        items={disorders()}
                        value={field().state.value}
                        onChange={(val) => field().handleChange(val)}
                        excludeIds={[d2, d3]}
                        renderSelected={renderDisorder}
                      />
                    )}
                  </props.form.Field>
                  <props.form.Field name="disorder2">
                    {(field) => (
                      <GlossaryAutocomplete
                        items={disorders()}
                        value={field().state.value}
                        onChange={(val) => field().handleChange(val)}
                        excludeIds={[d1, d3]}
                        renderSelected={renderDisorder}
                      />
                    )}
                  </props.form.Field>
                  <props.form.Field name="disorder3">
                    {(field) => (
                      <GlossaryAutocomplete
                        items={disorders()}
                        value={field().state.value}
                        onChange={(val) => field().handleChange(val)}
                        excludeIds={[d1, d2]}
                        renderSelected={renderDisorder}
                      />
                    )}
                  </props.form.Field>
                </div>
              </div>
            </details>
          );
        }}
      </props.form.Subscribe>
      <props.form.Subscribe selector={(state) => state.values.fightingArt}>
        {(fightingArt) => (
          <details class={styles.card}>
            <summary class="collapse-title">
              <span class={styles.cardTitle}>Fighting Art</span>
              {fightingArt() && <span class={styles.countBadge}>1</span>}
            </summary>
            <div class="collapse-content">
              <div class={styles.cardContent}>
                <props.form.Field name="fightingArt">
                  {(field) => (
                    <GlossaryAutocomplete
                      items={fightingArts()}
                      value={field().state.value}
                      onChange={(val) => field().handleChange(val)}
                      filter={(art) => !art.secret}
                      renderSelected={renderFightingArt}
                    />
                  )}
                </props.form.Field>
              </div>
            </div>
          </details>
        )}
      </props.form.Subscribe>
      <props.form.Subscribe selector={(state) => state.values.secretFightingArt}>
        {(secretFightingArt) => (
          <details class={styles.card}>
            <summary class="collapse-title">
              <span class={styles.cardTitle}>Secret Fighting Art</span>
              {secretFightingArt() && <span class={styles.countBadge}>1</span>}
            </summary>
            <div class="collapse-content">
              <div class={styles.cardContent}>
                <props.form.Field name="secretFightingArt">
                  {(field) => (
                    <GlossaryAutocomplete
                      items={fightingArts()}
                      value={field().state.value}
                      onChange={(val) => field().handleChange(val)}
                      filter={(art) => art.secret}
                      renderSelected={renderFightingArt}
                    />
                  )}
                </props.form.Field>
              </div>
            </div>
          </details>
        )}
      </props.form.Subscribe>
    </section>
  );
}
