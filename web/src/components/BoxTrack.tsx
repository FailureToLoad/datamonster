import styles from './BoxTrack.module.css';

interface BoxTrackProps {
  value: number;
  onChange: (val: number) => void;
  label: string;
  totalBoxes: number;
  accentedBoxes: number[];
  labelPosition?: 'top' | 'bottom' | 'left';
  class?: string;
}

export function BoxTrack(props: BoxTrackProps) {
  const handleBoxClick = (index: number) => {
    const boxNumber = index + 1;
    if (props.value === boxNumber) {
      props.onChange(boxNumber - 1);
    } else {
      props.onChange(boxNumber);
    }
  };

  const boxes = (
    <div class={styles.boxes}>
      {Array.from({length: props.totalBoxes}, (_, i) => {
        const boxNumber = i + 1;
        const isFilled = boxNumber <= props.value;
        const isAccented = props.accentedBoxes.includes(boxNumber);
        return (
          <button
            type="button"
            onClick={() => handleBoxClick(i)}
            class={`${styles.box} ${isAccented ? styles.boxAccented : ''} ${isFilled ? styles.boxFilled : ''}`}
          />
        );
      })}
    </div>
  );

  const labelElement = (
    <div class={styles.labelWrapper} data-tip={props.value}>
      <p class={styles.label}>{props.label}</p>
    </div>
  );

  const labelPosition = props.labelPosition ?? 'top';
  const containerClass = props.class ?? '';

  if (labelPosition === 'left') {
    return (
      <div class={`${styles.container} ${styles.containerRow} ${containerClass}`}>
        {labelElement}
        {boxes}
      </div>
    );
  }

  if (labelPosition === 'bottom') {
    return (
      <div class={`${styles.container} ${styles.containerColumn} ${containerClass}`}>
        {boxes}
        {labelElement}
      </div>
    );
  }

  return (
    <div class={`${styles.container} ${styles.containerColumn} ${containerClass}`}>
      {labelElement}
      {boxes}
    </div>
  );
}
