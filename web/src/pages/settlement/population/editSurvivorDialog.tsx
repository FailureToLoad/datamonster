import {useForm} from '@tanstack/react-form';
import {useRef, useLayoutEffect, useState} from 'react';
import {type Survivor, SurvivorGender, SurvivorStatus} from '~/types/survivor';
import {type} from 'arktype';
import {PatchJSON} from '~/lib/request';
import {BoxTrack} from '~/components/BoxTrack'
import { GenderFemaleIcon, GenderMaleIcon } from '@phosphor-icons/react';

const isInteger = type('number.integer');
const isPositive = type('number.integer >= 0');
const statusValidator = type.enumerated(
  SurvivorStatus.Alive, 
  SurvivorStatus.CannotDepart,
  SurvivorStatus.CeasedToExist,
  SurvivorStatus.Dead,
  SurvivorStatus.Retired,
);

const EditSurvivorSchema = type({
  status: statusValidator,
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
  birth: isPositive,
  huntxp: isPositive,
  courage: isPositive,
  understanding: isPositive,
});

type EditSurvivorFields = typeof EditSurvivorSchema.infer;

type EditSurvivorDialogProps = {
  survivor: Survivor | null;
  onClose: () => void;
  onSuccess: () => void;
};

export type SurvivorUpdateRequest = {
  statUpdates?: Record<string, number>;
  statusUpdate?: SurvivorStatus;
};

export default function EditSurvivorDialog({survivor, onClose, onSuccess}: EditSurvivorDialogProps) {
  const dialogRef = useRef<HTMLDialogElement>(null);
  const [cannotSpendSurvival, setCannotSpendSurvival] = useState(false);

  const form = useForm({
    defaultValues: {
      status: survivor?.status ?? SurvivorStatus.Alive,
      survival: survivor?.survival ?? 1,
      systemicPressure: survivor?.systemicPressure ?? 0,
      movement: survivor?.movement ?? 5,
      accuracy: survivor?.accuracy ?? 0,
      strength: survivor?.strength ?? 0,
      evasion: survivor?.evasion ?? 0,
      luck: survivor?.luck ?? 0,
      speed: survivor?.speed ?? 0,
      lumi: survivor?.lumi ?? 0,
      insanity: survivor?.insanity ?? 0,
      torment: survivor?.torment ?? 0,
      birth: survivor?.birth ?? 1,
      huntxp: survivor?.huntxp ?? 0,
      courage: survivor?.courage ?? 0,
      understanding: survivor?.understanding ?? 0,
    } as EditSurvivorFields,
    onSubmit: async ({value}) => {
      if (!survivor) return;

      const parsed = EditSurvivorSchema(value);
      if (parsed instanceof type.errors) {
        return;
      }

      const statFields: Array<[keyof EditSurvivorFields, keyof Survivor, string]> = [
        ['huntxp', 'huntxp', 'huntxp'],
        ['survival', 'survival', 'survival'],
        ['movement', 'movement', 'movement'],
        ['accuracy', 'accuracy', 'accuracy'],
        ['strength', 'strength', 'strength'],
        ['evasion', 'evasion', 'evasion'],
        ['luck', 'luck', 'luck'],
        ['speed', 'speed', 'speed'],
        ['lumi', 'lumi', 'lumi'],
        ['insanity', 'insanity', 'insanity'],
        ['torment', 'torment', 'torment'],
        ['systemicPressure', 'systemicPressure', 'systemicPressure'],
        ['courage', 'courage', 'courage'],
        ['understanding', 'understanding', 'understanding'],
      ];

      const statUpdates = statFields.reduce<Record<string, number>>((acc, [formKey, survivorKey, apiKey]) => {
        if (parsed[formKey] !== survivor[survivorKey]) {
          acc[apiKey] = parsed[formKey] as number;
        }
        return acc;
      }, {});

      const statusUpdate = parsed.status !== survivor.status ? parsed.status : undefined;

      if (Object.keys(statUpdates).length === 0 && !statusUpdate) {
        dialogRef.current?.close();
        onClose();
        return;
      }

      const payload: SurvivorUpdateRequest = {};
      if (Object.keys(statUpdates).length > 0) {
        payload.statUpdates = statUpdates;
      }
      if (statusUpdate) {
        payload.statusUpdate = statusUpdate;
      }

      const response = await PatchJSON(`/api/settlements/${survivor.settlementId}/survivors/${survivor.id}`, payload);

      if (response.ok) {
        dialogRef.current?.close();
        onSuccess();
      }
    },
  });

  useLayoutEffect(() => {
    if (survivor) {
      dialogRef.current?.showModal();
    }
  }, [survivor]);

  useLayoutEffect(() => {
    const dialog = dialogRef.current;
    if (!dialog) return;

    const handleDialogClose = () => {
      onClose();
    };

    dialog.addEventListener('close', handleDialogClose);
    return () => dialog.removeEventListener('close', handleDialogClose);
  }, [onClose]);

  const handleClose = () => {
    dialogRef.current?.close();
  };

  if (!survivor) return null;

  return (
    <dialog ref={dialogRef} className="modal">
      <div className="modal-box max-w-1/2 min-w-1/4 mx-auto px-6">
        <form
          onSubmit={(e) => {
            e.preventDefault();
            e.stopPropagation();
            form.handleSubmit();
          }}
        >
          <section className="grid grid-cols-2 items-center justify-center gap-4">
            <div id="base-details" className="col-span-2 flex flex-col">
              <div className="flex flex-row items-center justify-between border-b-2 border-black pb-1">
                <div className="flex flex-row items-center gap-2">
                  <p className="text-2xl">
                    {survivor.name}
                  </p>
                  <span className="badge badge-neutral">{survivor.gender === SurvivorGender.M ? < GenderMaleIcon weight="bold"  /> : <GenderFemaleIcon  weight="bold" />}</span>
                  <span className="badge badge-neutral font-semibold">Born year {survivor.birth}</span>
                </div>
                <form.Field
                  name="status"
                  validators={{
                    onChange: ({value}) => {
                      const result = statusValidator(value);
                      return result instanceof type.errors
                        ? result.summary
                        : undefined;
                    },
                  }}
                >
                  {(field) => (
                    <select
                      id="status-input"
                      className="select select-bordered select-sm w-auto min-w-0"
                      value={field.state.value}
                      onChange={(e) =>
                        field.handleChange(e.target.value as SurvivorStatus)
                      }
                    >
                      {Object.values(SurvivorStatus).map((status) => (
                        <option key={status} value={status}>
                          {status}
                        </option>
                      ))}
                    </select>
                  )}
                </form.Field>
              </div>
              <form.Field name="huntxp">
                {(field) => (
                  <HuntXPTrack
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                  />
                )}
              </form.Field>
              <div className="flex flex-row items-center gap-2 py-1">
                <button
                  type="button"
                  className={`badge font-semibold cursor-pointer ${cannotSpendSurvival ? 'badge-error' : 'badge-success'}`}
                  onClick={() => setCannotSpendSurvival(!cannotSpendSurvival)}
                >
                  {cannotSpendSurvival ? 'Cannot spend survival' : 'Can spend survival'}
                </button>
                <form.Subscribe selector={(state) => state.values.insanity}>
                  {(insanity) => insanity >= 3 && (
                    <span className="badge badge-neutral font-semibold">Insane</span>
                  )}
                </form.Subscribe>
              </div>
            </div>




            <div className="flex flex-row items-center justify-around col-span-2 border border-black gap-2">
              <form.Field name="survival">
                {(field) => (
                  <StatBox
                    id="survival-edit"
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                    label="Survival"
                    className="flex-1"
                  />
                )}
              </form.Field>
              <form.Field name="insanity">
                {(field) => (
                  <StatBox
                    id="insanity-edit"
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                    label="Insanity"
                    className="flex-1"
                  />
                )}
              </form.Field>
              <form.Field name="systemicPressure">
                {(field) => (
                  <StatBox
                    id="systemicPressure-edit"
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                    label="Systemic Pressure"
                    className="flex-1"
                  />
                )}
              </form.Field>
              <form.Field name="torment">
                {(field) => (
                  <StatBox
                    id="torment-edit"
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                    label="Torment"
                    className="flex-1"
                  />
                )}
              </form.Field>
              <form.Field name="lumi">
                {(field) => (
                  <StatBox
                    id="lumi-edit"
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                    label="Lumi"
                    className="flex-1"
                  />
                )}
              </form.Field>
            </div>

            <div id="base-stats" className="flex flex-row items-center justify-around col-span-2 border border-black gap-2">
              <form.Field name="movement">
                {(field) => (
                  <StatBox
                    id="movement-edit"
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                    label="Movement"
                    className="flex-1"
                  />
                )}
              </form.Field>
              <form.Field name="accuracy">
                {(field) => (
                  <StatBox
                    id="accuracy-edit"
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                    label="Accuracy"
                    className="flex-1"
                  />
                )}
              </form.Field>
              <form.Field name="strength">
                {(field) => (
                  <StatBox
                    id="strength-edit"
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                    label="Strength"
                    className="flex-1"
                  />
                )}
              </form.Field>
              <form.Field name="evasion">
                {(field) => (
                  <StatBox
                    id="evasion-edit"
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                    label="Evasion"
                    className="flex-1"
                  />
                )}
              </form.Field>
              <form.Field name="luck">
                {(field) => (
                  <StatBox
                    id="luck-edit"
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                    label="Luck"
                    className="flex-1"
                  />
                )}
              </form.Field>
              <form.Field name="speed">
                {(field) => (
                  <StatBox
                    id="speed-edit"
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                    label="Speed"
                    className="flex-1"
                  />
                )}
              </form.Field>
            </div>

            <div className="flex flex-row items-center col-span-2 border border-black gap-4 p-4">
              <form.Field name="courage">
                {(field) => (
                  <CourageTrack
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
                  />
                )}
              </form.Field>
              
              <form.Field name="understanding">
                {(field) => (
                  <UnderstandingTrack
                    value={field.state.value}
                    onChange={(val) => field.handleChange(val)}
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
              onClick={handleClose}
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
                  Save
                </button>
              )}
            </form.Subscribe>
          </div>
        </form>
      </div>
      <form method="dialog" className="modal-backdrop">
        <button onClick={handleClose}>close</button>
      </form>
    </dialog>
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
      className={`flex flex-col items-center ${className}`}
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
      <p className="my-4 text-xs text-center h-8 flex items-center">
        {label}
      </p>
    </div>
  );
}


function HuntXPTrack({
  value,
  onChange,
}: {
  value: number;
  onChange: (val: number) => void;
}) {
  return (
    <BoxTrack
      value={value}
      onChange={onChange}
      label="Hunt XP"
      totalBoxes={16}
      accentedBoxes={[2, 6, 10, 15, 16]}
      labelPosition="left"
    />
  );
}

function CourageTrack({
  value,
  onChange,
}: {
  value: number;
  onChange: (val: number) => void;
}) {
  return (
    <BoxTrack
      value={value}
      onChange={onChange}
      label="Courage"
      totalBoxes={9}
      accentedBoxes={[3, 9]}
      labelPosition="top"
      className="flex-1"
    />
  );
}

function UnderstandingTrack({
  value,
  onChange,
}: {
  value: number;
  onChange: (val: number) => void;
}) {
  return (
    <BoxTrack
      value={value}
      onChange={onChange}
      label="Understanding"
      totalBoxes={9}
      accentedBoxes={[3, 9]}
      labelPosition="top"
      className="flex-1"
    />
  );
}
