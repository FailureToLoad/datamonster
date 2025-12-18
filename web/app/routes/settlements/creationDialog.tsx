import * as z from 'zod';
import { useForm } from '@tanstack/react-form';
import { useRef } from 'react';
import { Plus } from 'lucide-react';

const AddSettlementSchema = z.object({
  settlementName: z
    .string()
    .min(1, 'Settlement name is required')
    .max(25, 'Settlement name is too long'),
});

type AddSettlementFields = z.infer<typeof AddSettlementSchema>;

export function CreateSettlementDialog({ refresh }: { refresh: () => void }) {
  const dialogRef = useRef<HTMLDialogElement>(null);

  const form = useForm({
    defaultValues: {
      settlementName: '',
    } as AddSettlementFields,
    onSubmit: async ({ value }) => {
      const parsed = AddSettlementSchema.safeParse(value);
      if (!parsed.success) {
        return;
      }

      const response = await fetch('/api/settlements', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: parsed.data.settlementName }),
      });

      if (response.ok) {
        refresh();
        dialogRef.current?.close();
        form.reset();
      }
    },
  });

  return (
    <>
      <button
        className="w-full btn btn-outline"
        aria-label="Create Settlement"
        onClick={() => dialogRef.current?.showModal()}
      >
        <Plus className="h-6 w-6" />
      </button>
      <dialog ref={dialogRef} className="modal">
        <div className="modal-box">
          <h3 className="font-bold text-lg">Add Settlement</h3>
          <form
            onSubmit={(e) => {
              e.preventDefault();
              e.stopPropagation();
              form.handleSubmit();
            }}
          >
            <form.Field
              name="settlementName"
              validators={{
                onChange: ({ value }) => {
                  const result = AddSettlementSchema.shape.settlementName.safeParse(value);
                  return result.success ? undefined : result.error.issues[0].message;
                },
              }}
            >
              {(field) => (
                <label className="form-control w-full">
                  <div className="label">
                    <span className="label-text">Settlement Name</span>
                  </div>
                  <input
                    type="text"
                    placeholder="Enter settlement name"
                    className="input input-bordered w-full"
                    value={field.state.value}
                    onChange={(e) => field.handleChange(e.target.value)}
                    onBlur={field.handleBlur}
                  />
                  {field.state.meta.errors.length > 0 && (
                    <div className="label">
                      <span className="label-text-alt text-error">
                        {field.state.meta.errors[0]}
                      </span>
                    </div>
                  )}
                </label>
              )}
            </form.Field>
            <div className="modal-action">
              <button
                type="button"
                className="btn"
                onClick={() => dialogRef.current?.close()}
              >
                Cancel
              </button>
              <form.Subscribe selector={(state) => state.canSubmit}>
                {(canSubmit) => (
                  <button type="submit" className="btn btn-primary" disabled={!canSubmit}>
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
