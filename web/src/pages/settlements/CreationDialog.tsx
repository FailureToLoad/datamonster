import { type } from "arktype";
import { useForm } from "@tanstack/react-form";
import { useRef } from "react";
import {PlusIcon} from "@phosphor-icons/react";

const settlementNameValidator = type("5 <= string <= 25");

const AddSettlementSchema = type({
  settlementName: settlementNameValidator,
});

type AddSettlementFields = typeof AddSettlementSchema.infer;

export function CreateSettlementDialog({ refresh }: { refresh: () => void }) {
  const dialogRef = useRef<HTMLDialogElement>(null);

  const form = useForm({
    defaultValues: {
      settlementName: "",
    } as AddSettlementFields,
    onSubmit: async ({ value }) => {
      const parsed = AddSettlementSchema(value);
      if (parsed instanceof type.errors) {
        return;
      }

      const response = await fetch("/api/settlements", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name: parsed.settlementName }),
        credentials: "include",
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
        className="w-96 btn btn-outline"
        aria-label="Create Settlement"
        onClick={() => dialogRef.current?.showModal()}
      >
        <PlusIcon className="h-6 w-6" />
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
                  const result = settlementNameValidator(value);
                  return result instanceof type.errors ? result.summary : undefined;
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
                  />
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
