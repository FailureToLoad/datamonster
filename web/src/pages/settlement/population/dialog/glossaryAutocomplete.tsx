import {useState, useRef} from 'react';
import styles from './glossaryAutocomplete.module.css';

type GlossaryItem = {
  id: string;
  name: string;
};

type GlossaryAutocompleteProps<T extends GlossaryItem> = {
  items: T[];
  value: string | null;
  onChange: (value: string | null) => void;
  excludeIds?: (string | null)[];
  filter?: (item: T) => boolean;
  renderSelected?: (item: T) => React.ReactNode;
  placeholder?: string;
};

export function GlossaryAutocomplete<T extends GlossaryItem>({
  items,
  value,
  onChange,
  excludeIds = [],
  filter,
  renderSelected,
  placeholder = 'None',
}: GlossaryAutocompleteProps<T>) {
  const [query, setQuery] = useState('');
  const [isOpen, setIsOpen] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const filteredByType = filter ? items.filter(filter) : items;
  const selectedItem = filteredByType.find((item) => item.id === value);
  const excludeSet = new Set(excludeIds.filter((id): id is string => id !== null && id !== value));

  const available = filteredByType.filter((item) => !excludeSet.has(item.id));
  const filtered = query
    ? available.filter((item) => item.name.toLowerCase().includes(query.toLowerCase()))
    : available;

  function handleSelect(item: T) {
    onChange(item.id);
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
    <div className={styles.container} data-dropdown-open={isOpen || undefined}>
      <div className={styles.inputWrapper}>
        {selectedItem ? (
          <div className={styles.selectedDisplay}>
            {renderSelected ? renderSelected(selectedItem) : selectedItem.name}
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
            placeholder={placeholder}
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
          {filtered.map((item) => (
            <li key={item.id}>
              <button
                type="button"
                className={styles.optionButton}
                onMouseDown={() => handleSelect(item)}
              >
                {item.name}
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
