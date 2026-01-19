import { createSignal, For, Show, type JSX } from 'solid-js';
import styles from './glossaryAutocomplete.module.css';

type GlossaryItem = {
  id: string;
  name: string;
};

interface GlossaryAutocompleteProps<T extends GlossaryItem> {
  items: T[];
  value: string | null;
  onChange: (value: string | null) => void;
  excludeIds?: (string | null)[];
  filter?: (item: T) => boolean;
  renderSelected?: (item: T) => JSX.Element;
  placeholder?: string;
}

export function GlossaryAutocomplete<T extends GlossaryItem>(props: GlossaryAutocompleteProps<T>) {
  const [query, setQuery] = createSignal('');
  const [isOpen, setIsOpen] = createSignal(false);
  let inputRef: HTMLInputElement | undefined;

  const filteredByType = () => (props.filter ? props.items.filter(props.filter) : props.items);
  const selectedItem = () => filteredByType().find((item) => item.id === props.value);
  const excludeSet = () => {
    const ids = props.excludeIds ?? [];
    return new Set(ids.filter((id): id is string => id !== null && id !== props.value));
  };

  const available = () => filteredByType().filter((item) => !excludeSet().has(item.id));
  const filtered = () => {
    const q = query();
    return q ? available().filter((item) => item.name.toLowerCase().includes(q.toLowerCase())) : available();
  };

  function handleSelect(item: T) {
    props.onChange(item.id);
    setQuery('');
    setIsOpen(false);
  }

  function handleClear() {
    props.onChange(null);
    setQuery('');
    inputRef?.focus();
  }

  function handleInputChange(e: Event) {
    const target = e.target as HTMLInputElement;
    setQuery(target.value);
    if (props.value) props.onChange(null);
    setIsOpen(true);
  }

  const placeholder = () => props.placeholder ?? 'None';

  return (
    <div class={styles.container} data-dropdown-open={isOpen() || undefined}>
      <div class={styles.inputWrapper}>
        <Show
          when={selectedItem()}
          fallback={
            <input
              ref={inputRef}
              type="text"
              class={styles.input}
              value={query()}
              onInput={handleInputChange}
              onFocus={() => setIsOpen(true)}
              onBlur={() => setTimeout(() => setIsOpen(false), 150)}
              placeholder={placeholder()}
            />
          }
        >
          {(item) => (
            <div class={styles.selectedDisplay}>
              {props.renderSelected ? props.renderSelected(item()) : item().name}
            </div>
          )}
        </Show>
        <Show when={props.value}>
          <button type="button" class={styles.clearButton} onClick={handleClear}>
            ×
          </button>
        </Show>
      </div>
      <Show when={isOpen() && filtered().length > 0}>
        <ul class={styles.dropdownList}>
          <For each={filtered()}>
            {(item) => (
              <li>
                <button type="button" class={styles.optionButton} onMouseDown={() => handleSelect(item)}>
                  {item.name}
                </button>
              </li>
            )}
          </For>
        </ul>
      </Show>
    </div>
  );
}
