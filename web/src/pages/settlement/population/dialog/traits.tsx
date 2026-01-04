import {withForm} from '~/lib/form';
import type {SurvivorFormFields} from './form';
import styles from './traits.module.css';
import { DisorderAutocomplete } from './disorderAutocomplete';

export const TraitsTab = withForm({
  defaultValues: {} as SurvivorFormFields,
  render: function Render({form}) {
    return (
      <section className={styles.section}>
        <div className={styles.disordersCard}>
          <h3 className={styles.cardTitle}>Disorders</h3>
          <div className={styles.disordersList}>
            <form.Subscribe selector={(state) => [state.values.disorder1, state.values.disorder2, state.values.disorder3]}>
              {([d1, d2, d3]) => (
                <>
                  <form.Field name="disorder1">
                    {(field) => (
                      <DisorderAutocomplete
                        value={field.state.value}
                        onChange={(val) => field.handleChange(val)}
                        excludeIds={[d2, d3]}
                      />
                    )}
                  </form.Field>
                  <form.Field name="disorder2">
                    {(field) => (
                      <DisorderAutocomplete
                        value={field.state.value}
                        onChange={(val) => field.handleChange(val)}
                        excludeIds={[d1, d3]}
                      />
                    )}
                  </form.Field>
                  <form.Field name="disorder3">
                    {(field) => (
                      <DisorderAutocomplete
                        value={field.state.value}
                        onChange={(val) => field.handleChange(val)}
                        excludeIds={[d1, d2]}
                      />
                    )}
                  </form.Field>
                </>
              )}
            </form.Subscribe>
          </div>
        </div>
      </section>
    );
  },
});
