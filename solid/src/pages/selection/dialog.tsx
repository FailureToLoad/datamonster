import { type } from "arktype";
import { createForm } from "@tanstack/solid-form";
import Plus from "lucide-solid/icons/plus";
import { PostJSON } from "~/lib/request";
import styles from "./dialog.module.css";

const settlementNameValidator = type("5 <= string <= 25");

const AddSettlementSchema = type({
  settlementName: settlementNameValidator,
});

type AddSettlementFields = typeof AddSettlementSchema.infer;

export function CreateSettlementDialog(props: { refresh: () => void }) {
  let dialogRef: HTMLDialogElement | undefined;

  const form = createForm(() => ({
    defaultValues: {
      settlementName: "",
    } as AddSettlementFields,
    onSubmit: async ({ value }) => {
      const parsed = AddSettlementSchema(value);
      if (parsed instanceof type.errors) {
        return;
      }

      const response = await PostJSON("/api/settlements", { name: parsed.settlementName });

      if (response.ok) {
        props.refresh();
        dialogRef?.close();
        form.reset();
      }
    },
  }));

  return (
    <>
      <button
        class={`btn btn-neutral ${styles.createButton}`}
        aria-label="Create Settlement"
        onClick={() => dialogRef?.showModal()}
      >
        <Plus class={styles.icon} />
      </button>
      <dialog ref={dialogRef} class="modal">
        <div class="modal-box">
          <h3 class={styles.dialogTitle}>Add Settlement</h3>
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
                <input
                  type="text"
                  placeholder="Enter settlement name"
                  class={`input ${styles.formInput}`}
                  value={field().state.value}
                  onInput={(e) => field().handleChange(e.currentTarget.value)}
                />
              )}
            </form.Field>
            <div class="modal-action">
              <button
                type="button"
                class="btn"
                onClick={() => dialogRef?.close()}
              >
                Cancel
              </button>
              <form.Subscribe selector={(state) => state.canSubmit}>
                {(canSubmit) => (
                  <button type="submit" class="btn btn-primary" disabled={!canSubmit()}>
                    Create
                  </button>
                )}
              </form.Subscribe>
            </div>
          </form>
        </div>
        <form method="dialog" class="modal-backdrop">
          <button>close</button>
        </form>
      </dialog>
    </>
  );
}
