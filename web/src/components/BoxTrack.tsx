import styles from './BoxTrack.module.css';

export function BoxTrack({
  value,
  onChange,
  label,
  totalBoxes,
  accentedBoxes,
  labelPosition = 'top',
  className = '',
}: {
  value: number;
  onChange: (val: number) => void;
  label: string;
  totalBoxes: number;
  accentedBoxes: number[];
  labelPosition?: 'top' | 'bottom' | 'left';
  className?: string;
}) {
  const handleBoxClick = (index: number) => {
    const boxNumber = index + 1;
    if (value === boxNumber) {
      onChange(boxNumber - 1);
    } else {
      onChange(boxNumber);
    }
  };

  const boxes = (
    <div className={styles.boxes}>
      {Array.from({length: totalBoxes}, (_, i) => {
        const boxNumber = i + 1;
        const isFilled = boxNumber <= value;
        const isAccented = accentedBoxes.includes(boxNumber);
        return (
          <button
            key={i}
            type="button"
            onClick={() => handleBoxClick(i)}
            className={`${styles.box} ${isAccented ? styles.boxAccented : ''} ${isFilled ? styles.boxFilled : ''}`}
          />
        );
      })}
    </div>
  );

  const labelElement = (
    <div className={styles.labelWrapper} data-tip={value}>
      <p className={styles.label}>{label}</p>
    </div>
  );

  if (labelPosition === 'left') {
    return (
      <div className={`${styles.container} ${styles.containerRow} ${className}`}>
        {labelElement}
        {boxes}
      </div>
    );
  }

  if (labelPosition === 'bottom') {
    return (
      <div className={`${styles.container} ${styles.containerColumn} ${className}`}>
        {boxes}
        {labelElement}
      </div>
    );
  }

  return (
    <div className={`${styles.container} ${styles.containerColumn} ${className}`}>
      {labelElement}
      {boxes}
    </div>
  );
}