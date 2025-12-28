
export function BoxTrack({
  value,
  onChange,
  label,
  totalBoxes,
  accentedBoxes,
  labelPosition = 'top',
}: {
  value: number;
  onChange: (val: number) => void;
  label: string;
  totalBoxes: number;
  accentedBoxes: number[];
  labelPosition?: 'top' | 'bottom' | 'left';
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
    <div className="flex flex-row flex-wrap gap-2">
      {Array.from({length: totalBoxes}, (_, i) => {
        const boxNumber = i + 1;
        const isFilled = boxNumber <= value;
        const isAccented = accentedBoxes.includes(boxNumber);
        return (
          <button
            key={i}
            type="button"
            onClick={() => handleBoxClick(i)}
            className={`size-5 border transition-colors ${
              isAccented ? 'border-2 border-black' : 'border border-gray-400'
            } ${isFilled ? 'bg-black' : 'bg-white hover:bg-gray-200'}`}
          />
        );
      })}
    </div>
  );

  const labelElement = (
    <div className="tooltip" data-tip={value}>
      <p className="text-sm cursor-help">{label}</p>
    </div>
  );

  if (labelPosition === 'left') {
    return (
      <div className="flex flex-row items-center gap-2 py-2">
        {labelElement}
        {boxes}
      </div>
    );
  }

  if (labelPosition === 'bottom') {
    return (
      <div className="flex flex-col items-start gap-2 py-2">
        {boxes}
        {labelElement}
      </div>
    );
  }

  return (
    <div className="flex flex-col items-start gap-2 py-2">
      {labelElement}
      {boxes}
    </div>
  );
}