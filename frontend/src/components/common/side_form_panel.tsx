import type { ReactNode, FormEvent } from "react";

interface SideFormPanelProps {
  open: boolean;
  onClose: () => void;
  title: string;
  onSubmit: (e: FormEvent) => void;
  submitting: boolean;
  submitLabel: string;
  savingLabel?: string;
  /** Tailwind width class for the panel, e.g. "w-80" (default) or "w-96". */
  widthClassName?: string;
  /** Set true when the form has many fields and needs its own scroll (e.g. project_task_form_panel). */
  scrollable?: boolean;
  children: ReactNode;
}

/**
 * Shared shell for every "New/Edit X" slide-over panel: backdrop, right-side
 * panel, title + close button, and a Cancel/Save footer pinned to the bottom.
 * Individual forms only need to render their fields as children.
 */
export default function SideFormPanel({
  open,
  onClose,
  title,
  onSubmit,
  submitting,
  submitLabel,
  savingLabel = "Saving...",
  widthClassName = "w-80",
  scrollable = false,
  children,
}: SideFormPanelProps) {
  if (!open) return null;

  return (
    <>
      <div
        className="fixed inset-0 z-10 bg-black/20 transition-opacity"
        onClick={onClose}
      />

      <div
        className={`fixed right-0 top-0 z-20 flex h-full ${widthClassName} flex-col ${
          scrollable ? "overflow-y-auto" : ""
        } border-l border-(--color-border) bg-(--color-surface) p-6 shadow-lg transition-transform`}
      >
        <div className="mb-6 flex items-center justify-between">
          <h2 className="app-panel-title text-lg">{title}</h2>
          <button
            onClick={onClose}
            className="text-(--color-muted) hover:text-(--color-text)"
          >
            <i className="bx bx-x text-2xl"></i>
          </button>
        </div>

        <form onSubmit={onSubmit} className="flex h-full flex-col space-y-4">
          {children}

          <div className="mt-auto flex justify-end gap-2 pt-6">
            <button
              type="button"
              onClick={onClose}
              className="app-btn-secondary px-4 py-2 text-sm font-medium"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={submitting}
              className="app-btn-primary px-4 py-2 text-sm font-medium disabled:opacity-50"
            >
              {submitting ? savingLabel : submitLabel}
            </button>
          </div>
        </form>
      </div>
    </>
  );
}