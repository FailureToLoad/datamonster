import {withForm} from '~/lib/form';
import type {SurvivorFormFields} from './form';
import styles from './traits.module.css';
import {GlossaryAutocomplete} from './glossaryAutocomplete';
import {useGlossary} from '~/hooks/glossary';
import type {Disorder, FightingArt} from '~/lib/glossary';

function renderDisorder(disorder: Disorder) {
  return (
    <>
      <span className={styles.itemName}>{disorder.name}</span>
      <span className={styles.itemEffect}> - {disorder.effect}</span>
    </>
  );
}

function renderFightingArt(art: FightingArt) {
  return (
    <>
      <span className={styles.itemName}>{art.name}</span>
      <span className={styles.itemEffect}> - {art.text.join(' ')}</span>
    </>
  );
}

export const TraitsTab = withForm({
  defaultValues: {} as SurvivorFormFields,
  render: function Render({form}) {
    const {disorders, fightingArts} = useGlossary();

    return (
      <section className={styles.section}>
        <form.Subscribe selector={(state) => [state.values.disorder1, state.values.disorder2, state.values.disorder3]}>
          {([d1, d2, d3]) => {
            const count = [d1, d2, d3].filter(Boolean).length;
            return (
              <details className={styles.card}>
                <summary className="collapse-title">
                  <span className={styles.cardTitle}>Disorders</span>
                  {count > 0 && <span className={styles.countBadge}>{count}</span>}
                </summary>
                <div className="collapse-content">
                  <div className={styles.cardContent}>
                    <form.Field name="disorder1">
                      {(field) => (
                        <GlossaryAutocomplete
                          items={disorders}
                          value={field.state.value}
                          onChange={(val) => field.handleChange(val)}
                          excludeIds={[d2, d3]}
                          renderSelected={renderDisorder}
                        />
                      )}
                    </form.Field>
                    <form.Field name="disorder2">
                      {(field) => (
                        <GlossaryAutocomplete
                          items={disorders}
                          value={field.state.value}
                          onChange={(val) => field.handleChange(val)}
                          excludeIds={[d1, d3]}
                          renderSelected={renderDisorder}
                        />
                      )}
                    </form.Field>
                    <form.Field name="disorder3">
                      {(field) => (
                        <GlossaryAutocomplete
                          items={disorders}
                          value={field.state.value}
                          onChange={(val) => field.handleChange(val)}
                          excludeIds={[d1, d2]}
                          renderSelected={renderDisorder}
                        />
                      )}
                    </form.Field>
                  </div>
                </div>
              </details>
            );
          }}
        </form.Subscribe>
        <form.Subscribe selector={(state) => state.values.fightingArt}>
          {(fightingArt) => (
            <details className={styles.card}>
              <summary className="collapse-title">
                <span className={styles.cardTitle}>Fighting Art</span>
                {fightingArt && <span className={styles.countBadge}>1</span>}
              </summary>
              <div className="collapse-content">
                <div className={styles.cardContent}>
                  <form.Field name="fightingArt">
                    {(field) => (
                      <GlossaryAutocomplete
                        items={fightingArts}
                        value={field.state.value}
                        onChange={(val) => field.handleChange(val)}
                        filter={(art) => !art.secret}
                        renderSelected={renderFightingArt}
                      />
                    )}
                  </form.Field>
                </div>
              </div>
            </details>
          )}
        </form.Subscribe>
        <form.Subscribe selector={(state) => state.values.secretFightingArt}>
          {(secretFightingArt) => (
            <details className={styles.card}>
              <summary className="collapse-title">
                <span className={styles.cardTitle}>Secret Fighting Art</span>
                {secretFightingArt && <span className={styles.countBadge}>1</span>}
              </summary>
              <div className="collapse-content">
                <div className={styles.cardContent}>
                  <form.Field name="secretFightingArt">
                    {(field) => (
                      <GlossaryAutocomplete
                        items={fightingArts}
                        value={field.state.value}
                        onChange={(val) => field.handleChange(val)}
                        filter={(art) => art.secret}
                        renderSelected={renderFightingArt}
                      />
                    )}
                  </form.Field>
                </div>
              </div>
            </details>
          )}
        </form.Subscribe>
      </section>
    );
  },
});
