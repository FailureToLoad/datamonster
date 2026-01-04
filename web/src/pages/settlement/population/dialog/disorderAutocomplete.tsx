import {useState, useRef} from 'react';
import {useGlossary} from '~/hooks/glossary';
import type {Disorder} from '~/lib/glossary';
import styles from './disorderAutocomplete.module.css';

type DisorderAutocompleteProps = {
  value: string | null;
  onChange: (value: string | null) => void;
  excludeIds?: (string | null)[];
};

export function DisorderAutocomplete({value, onChange, excludeIds = []}: DisorderAutocompleteProps) {
  const {disorders} = useGlossary();
  const [query, setQuery] = useState('');
  const [isOpen, setIsOpen] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const selectedDisorder = disorders.find((d) => d.id === value);
  const excludeSet = new Set(excludeIds.filter((id): id is string => id !== null && id !== value));

  const available = disorders.filter((d) => !excludeSet.has(d.id));
  const filtered = query
    ? available.filter((d) => d.name.toLowerCase().includes(query.toLowerCase()))
    : available;

  function handleSelect(disorder: Disorder) {
    onChange(disorder.id);
    setQuery('');
    setIsOpen(false);
  }

  function handleClear() {
    onChange(null);
    setQuery('');
    inputRef.current?.focus();
  }

  function handleInputChange(e: React.ChangeEvent<HTMLInputElement>) {
    setQuery(e.target.value);
    if (value) onChange(null);
    setIsOpen(true);
  }

  return (
    <div className={styles.container}>
      <div className={styles.inputWrapper}>
        {selectedDisorder ? (
          <div className={styles.selectedDisplay}>
            <span className={styles.selectedName}>{selectedDisorder.name}</span>
            <span className={styles.selectedEffect}> - {selectedDisorder.effect}</span>
          </div>
        ) : (
          <input
            ref={inputRef}
            type="text"
            className={styles.input}
            value={query}
            onChange={handleInputChange}
            onFocus={() => setIsOpen(true)}
            onBlur={() => setTimeout(() => setIsOpen(false), 150)}
            placeholder="None"
          />
        )}
        {value && (
          <button
            type="button"
            className={styles.clearButton}
            onClick={handleClear}
          >
            ×
          </button>
        )}
      </div>
      {isOpen && filtered.length > 0 && (
        <ul className={styles.dropdownList}>
          {filtered.map((disorder) => (
            <li key={disorder.id}>
              <button
                type="button"
                className={styles.optionButton}
                onMouseDown={() => handleSelect(disorder)}
              >
                {disorder.name}
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}