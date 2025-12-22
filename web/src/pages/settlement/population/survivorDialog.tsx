import {useForm} from '@tanstack/react-form';
import {useRef} from 'react';
import {PlusIcon} from '@phosphor-icons/react';
import {SurvivorGender} from '~/types/survivor';
import {type} from 'arktype';

const survivorNameValidator = type('1 <= string <= 50');
const isInteger = type('number.integer');
const isPositive = type('number.integer >= 0');
const genderValidator = type.enumerated(SurvivorGender.M, SurvivorGender.F);

const AddSurvivorSchema = type({
  name: survivorNameValidator,
  gender: genderValidator,
  survival: isPositive,
  systemicPressure: isInteger,
  movement: isInteger,
  accuracy: isInteger,
  strength: isInteger,
  evasion: isInteger,
  luck: isInteger,
  speed: isInteger,
  lumi: isPositive,
  insanity: isPositive,
  torment: isInteger,
});

type AddSurvivorFields = typeof AddSurvivorSchema.infer;

type SurvivorDialogProps = {
  settlementId: string;
  onSuccess?: () => void;
};

export default function NewSurvivorDialog({settlementId, onSuccess}: SurvivorDialogProps) {
  const dialogRef = useRef<HTMLDialogElement>(null);

  const form = useForm({
    defaultValues: {
      name: 'Meat',
      gender: SurvivorGender.M,
      survival: 1,
      systemicPressure: 0,
      movement: 5,
      accuracy: 0,
      strength: 0,
      evasion: 0,
      luck: 0,
      speed: 0,
      lumi: 0,
      insanity: 0,
      torment: 0,
    } as AddSurvivorFields,
    onSubmit: async ({value}) => {
      const parsed = AddSurvivorSchema(value);
      if (parsed instanceof type.errors) {
        return;
      }

      const response = await fetch(`/api/settlements/${settlementId}/survivors`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
          name: parsed.name,
          gender: parsed.gender,
          survival: parsed.survival,
          systemicpressure: parsed.systemicPressure,
          movement: parsed.movement,
          accuracy: parsed.accuracy,
          strength: parsed.strength,
          evasion: parsed.evasion,
          luck: parsed.luck,
          speed: parsed.speed,
          lumi: parsed.lumi,
          insanity: parsed.insanity,
          torment: parsed.torment,
          born: 1,
          huntxp: 0,
          courage: 0,
          understanding: 0,
          settlementID: settlementId,
        }),
        credentials: 'include',
      });

      if (response.ok) {
        dialogRef.current?.close();
        form.reset();
        onSuccess?.();
      }
    },
  });

  return (
    <>
      <button
        className="btn btn-outline"
        aria-label="Create Survivor"
        onClick={() => dialogRef.current?.showModal()}
      >
        <PlusIcon className="size-4" />
      </button>
      <dialog ref={dialogRef} className="modal mx-auto w-3/5 px-6 grow">
        <div className="modal-box w-3/5 max-w-none px-6">
          <form
            onSubmit={(e) => {
              e.preventDefault();
              e.stopPropagation();
              form.handleSubmit();
            }}
          >
            <section className="grid grid-cols-2 items-center justify-center gap-4">
              {/* Section 1: Name + Gender row */}
              <div className="mb-4 flex flex-row items-center justify-between col-span-2 h-full border-b-2 border-black">
                <div className="flex flex-row gap-2 w-fill">
                  <p className="text-2xl font-serif font-light tracking-wide">
                    Name
                  </p>
                  <form.Field
                    name="name"
                    validators={{
                      onChange: ({value}) => {
                        const result = survivorNameValidator(value);
                        return result instanceof type.errors
                          ? result.summary
                          : undefined;
                      },
                    }}
                  >
                    {(field) => (
                      <input
                        type="text"
                        id="name-input"
                        className="input input-bordered w-full text-lg"
                        value={field.state.value}
                        onChange={(e) => field.handleChange(e.target.value)}
                      />
                    )}
                  </form.Field>
                </div>
                <form.Field name="gender">
                  {(field) => (
                    <div className="flex flex-row gap-4">
                      <label className="flex items-center gap-2 cursor-pointer">
                        <input
                          type="radio"
                          name="gender"
                          className="radio"
                          value="M"
                          checked={field.state.value === 'M'}
                          onChange={() => field.handleChange(SurvivorGender.M)}
                        />
                        <span>M</span>
                      </label>
                      <label className="flex items-center gap-2 cursor-pointer">
                        <input
                          type="radio"
                          name="gender"
                          className="radio"
                          value="F"
                          checked={field.state.value === 'F'}
                          onChange={() => field.handleChange(SurvivorGender.F)}
                        />
                        <span>F</span>
                      </label>
                    </div>
                  )}
                </form.Field>
              </div>

              {/* Section 2: Survival box */}
              <div className="flex flex-row items-center col-span-2 border border-black h-full gap-2">
                <form.Field
                  name="survival"
                  validators={{
                    onChange: ({value}) => {
                      const result = isPositive(value);
                      return result instanceof type.errors
                        ? result.summary
                        : undefined;
                    },
                  }}
                >
                  {(field) => (
                    <div className="order-1 ml-6 my-4 border border-black size-20 place-content-around">
                      <input
                        id="survival-input"
                        type="number"
                        className="input input-bordered size-full text-center text-2xl border-0"
                        value={field.state.value}
                        onChange={(e) =>
                          field.handleChange(Number(e.target.value))
                        }
                      />
                    </div>
                  )}
                </form.Field>

                <div className="order-2 my-4 flex flex-col flex-1 h-20 items-start justify-between">
                  <p className="text-2xl font-serif font-light tracking-wide">
                    Survival
                  </p>
                  <div className="flex w-full">
                    <div className="flex items-start space-x-2">
                      <input
                        type="checkbox"
                        id="cannot-spend"
                        className="checkbox"
                      />
                      <label
                        htmlFor="cannot-spend"
                        className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                      >
                        Cannot spend survival
                      </label>
                    </div>
                  </div>
                </div>
                <div className="order-3 my-4 items-start justify-start content-start space-y-1">
                  <CheckboxItem id="dodge" label="Dodge" />
                  <CheckboxItem id="encourage" label="Encourage" />
                  <CheckboxItem id="surge" label="Surge" />
                  <CheckboxItem id="dash" label="Dash" />
                  <CheckboxItem id="fistpump" label="Fist Pump" />
                </div>
                <div className="order-last border-l border-black justify-self-end flex flex-col">
                  <form.Field name="systemicPressure">
                    {(field) => (
                      <StatBox
                        id="systemicPressure"
                        value={field.state.value}
                        onChange={(val) => field.handleChange(val)}
                        label="Systemic Pressure"
                        className="mx-4 size-min"
                      />
                    )}
                  </form.Field>
                </div>
              </div>

              {/* Section 3: Stats row */}
              <div className="flex flex-row items-center justify-between col-span-2 border border-black h-32 gap-2">
                <form.Field name="movement">
                  {(field) => (
                    <StatBox
                      id="movement"
                      value={field.state.value}
                      onChange={(val) => field.handleChange(val)}
                      label="Movement"
                      className="ml-4 mr-2 size-full"
                    />
                  )}
                </form.Field>
                <div className="divider divider-horizontal bg-black m-0 w-px" />
                <form.Field name="accuracy">
                  {(field) => (
                    <StatBox
                      id="accuracy"
                      value={field.state.value}
                      onChange={(val) => field.handleChange(val)}
                      label="Accuracy"
                      className="mx-2 size-full"
                    />
                  )}
                </form.Field>
                <div className="divider divider-horizontal bg-black m-0 w-px" />
                <form.Field name="strength">
                  {(field) => (
                    <StatBox
                      id="strength"
                      value={field.state.value}
                      onChange={(val) => field.handleChange(val)}
                      label="Strength"
                      className="mx-2 size-full"
                    />
                  )}
                </form.Field>
                <div className="divider divider-horizontal bg-black m-0 w-px" />
                <form.Field name="evasion">
                  {(field) => (
                    <StatBox
                      id="evasion"
                      value={field.state.value}
                      onChange={(val) => field.handleChange(val)}
                      label="Evasion"
                      className="mx-2 size-full"
                    />
                  )}
                </form.Field>
                <div className="divider divider-horizontal bg-black m-0 w-px" />
                <form.Field name="luck">
                  {(field) => (
                    <StatBox
                      id="luck"
                      value={field.state.value}
                      onChange={(val) => field.handleChange(val)}
                      label="Luck"
                      className="mx-2 size-full"
                    />
                  )}
                </form.Field>
                <div className="divider divider-horizontal bg-black m-0 w-px" />
                <form.Field name="speed">
                  {(field) => (
                    <StatBox
                      id="speed"
                      value={field.state.value}
                      onChange={(val) => field.handleChange(val)}
                      label="Speed"
                      className="mx-2 size-full"
                    />
                  )}
                </form.Field>
                <div className="divider divider-horizontal bg-black m-0 w-px" />
                <form.Field name="lumi">
                  {(field) => (
                    <StatBox
                      id="lumi"
                      value={field.state.value}
                      onChange={(val) => field.handleChange(val)}
                      label="Lumi"
                      className="ml-2 mr-4 size-full"
                    />
                  )}
                </form.Field>
              </div>

              {/* Section 4: Brain row */}
              <div className="flex flex-row items-center col-span-2 border border-black h-full gap-2">
                <form.Field name="insanity">
                  {(field) => (
                    <StatBox
                      id="insanity"
                      value={field.state.value}
                      onChange={(val) => field.handleChange(val)}
                      label="Insanity"
                      className="order-first ml-4 mr-2"
                    />
                  )}
                </form.Field>
                <div className="divider divider-horizontal bg-black m-0 w-px" />
                <div className="my-4 flex flex-col flex-1 h-20 items-start justify-between">
                  <div className="w-full flex flex-row justify-between">
                    <p className="text-2xl font-serif font-light tracking-wide">
                      Brain
                    </p>
                    <input
                      type="checkbox"
                      id="brainbox"
                      className="checkbox size-6"
                    />
                  </div>

                  <div className="flex w-full">
                    <p className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                      If your insanity is 3+, you are <b>insane</b>
                    </p>
                  </div>
                </div>
                <div className="divider divider-horizontal bg-black m-0 w-px" />
                <form.Field name="torment">
                  {(field) => (
                    <StatBox
                      id="torment"
                      value={field.state.value}
                      onChange={(val) => field.handleChange(val)}
                      label="Torment"
                      className="order-last ml-2 mr-4"
                    />
                  )}
                </form.Field>
              </div>
            </section>

            {/* Footer */}
            <div className="modal-action pt-4">
              <button
                type="button"
                className="btn"
                onClick={() => dialogRef.current?.close()}
              >
                Cancel
              </button>
              <form.Subscribe selector={(state) => state.canSubmit}>
                {(canSubmit) => (
                  <button
                    type="submit"
                    className="btn btn-primary"
                    disabled={!canSubmit}
                  >
                    Create
                  </button>
                )}
              </form.Subscribe>
            </div>
          </form>
        </div>
        <form method="dialog" className="modal-backdrop">
          <button>close</button>
        </form>
      </dialog>
    </>
  );
}

function CheckboxItem({id, label}: {id: string; label: string}) {
  return (
    <div className="flex items-start space-x-1 rounded-none">
      <input type="checkbox" id={id} className="checkbox checkbox-sm" />
      <label
        htmlFor={id}
        className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
      >
        {label}
      </label>
    </div>
  );
}

function StatBox({
  id,
  value,
  onChange,
  label,
  className = '',
}: {
  id: string;
  value: number;
  onChange: (val: number) => void;
  label: string;
  className?: string;
}) {
  return (
    <div
      className={`flex flex-col items-center justify-between w-fit ${className}`}
    >
      <div className="mt-4 flex border border-black h-16 w-14">
        <input
          id={`${id}-input`}
          type="number"
          className="input input-bordered size-full text-center text-lg border-0"
          value={value}
          onChange={(e) => onChange(Number(e.target.value))}
        />
      </div>
      <p className="flex my-4 text-xs text-wrap whitespace-break-spaces text-center">
        {label}
      </p>
    </div>
  );
}
