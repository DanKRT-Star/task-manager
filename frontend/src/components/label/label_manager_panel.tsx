import { useState } from "react";
import type { Label, CreateLabelPayload } from "../../types/label";

interface LabelManagerPanelProps {
  labels: Label[];
  loading: boolean;
  onCreate: (payload: CreateLabelPayload) => Promise<void>;
  onDelete: (labelId: number) => Promise<void>;
}

const DEFAULT_COLOR = "#6b7280";

export default function LabelManagerPanel({
  labels,
  loading,
  onCreate,
  onDelete,
}: LabelManagerPanelProps) {
  const [name, setName] = useState("");
  const [color, setColor] = useState(DEFAULT_COLOR);
  const [submitting, setSubmitting] = useState(false);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    setSubmitting(true);
    try {
      await onCreate({ name: name.trim(), color });
      setName("");
      setColor(DEFAULT_COLOR);
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = async (labelId: number) => {
    if (!confirm("Delete this label? It will be removed from every task.")) return;
    await onDelete(labelId);
  };

  return (
    <div className="app-card p-5">
      <form onSubmit={handleCreate} className="mb-5 flex flex-wrap items-center gap-2">
        <input
          type="color"
          value={color}
          onChange={(e) => setColor(e.target.value)}
          className="h-9 w-9 shrink-0 cursor-pointer rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) p-1"
          title="Label color"
        />
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Label name, e.g. bug"
          className="flex-1 rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
        />
        <button
          type="submit"
          disabled={submitting}
          className="app-btn-primary flex items-center gap-1 px-4 py-2 text-sm font-medium disabled:opacity-50"
        >
          <i className="bx bx-plus"></i>
          Add
        </button>
      </form>

      {loading ? (
        <p className="app-text-muted py-6 text-center text-sm">Loading labels...</p>
      ) : labels.length === 0 ? (
        <p className="app-text-muted py-6 text-center text-sm">No labels yet.</p>
      ) : (
        <div className="flex flex-wrap gap-2">
          {labels.map((label) => (
            <span
              key={label.labelId}
              className="flex items-center gap-1.5 rounded-full px-3 py-1 text-xs font-medium text-white"
              style={{ backgroundColor: label.color }}
            >
              {label.name}
              <button
                onClick={() => handleDelete(label.labelId)}
                className="ml-1 opacity-80 hover:opacity-100"
                title="Delete label"
              >
                <i className="bx bx-x text-sm"></i>
              </button>
            </span>
          ))}
        </div>
      )}
    </div>
  );
}