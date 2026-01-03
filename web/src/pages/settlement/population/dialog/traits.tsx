import {withForm} from '~/lib/form';
import type {SurvivorFormFields} from './form';
import styles from './traits.module.css';

export const TraitsTab = withForm({
  defaultValues: {} as SurvivorFormFields,
  render: function Render() {
    return (
      <section className={styles.section}>
        Disorders and Fighting Arts
      </section>
    );
  },
});
