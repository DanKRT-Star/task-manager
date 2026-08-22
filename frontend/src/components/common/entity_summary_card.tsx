import type { ReactNode } from "react";

interface EntitySummaryCardProps {
  title: string;
  /** Omit entirely (leave undefined) if the entity has no description field, e.g. Sprint. */
  description?: string;
  /** Optional badge rendered under the title, e.g. Sprint's status pill. */
  badge?: ReactNode;
  /** Optional meta line rendered under the description, e.g. an icon + due date. */
  meta?: ReactNode;
  onEdit: () => void;
  onDelete: () => void;
  onViewTasks: () => void;
  viewTasksLabel: string;
}

/**
 * Shared card used for Epic / Milestone / Sprint list tabs on the project
 * detail page: title + edit/delete icons, optional badge, description,
 * optional meta line, and a "view tasks" link.
 */
export default function EntitySummaryCard({
  title,
  description,
  badge,
  meta,
  onEdit,
  onDelete,
  onViewTasks,
  viewTasksLabel,
}: EntitySummaryCardProps) {
  return (
    <div className="app-card p-4">
      <div className="mb-2 flex items-start justify-between gap-3">
        <h3 className="text-sm font-bold text-(--color-text)">{title}</h3>
        <div className="flex shrink-0 gap-2">
          <button onClick={onEdit} className="text-blue-500 hover:text-blue-700">
            <i className="bx bx-edit text-base"></i>
          </button>
          <button onClick={onDelete} className="text-red-500 hover:text-red-700">
            <i className="bx bx-trash text-base"></i>
          </button>
        </div>
      </div>

      {badge}

      {description !== undefined && (
        <p className="app-text-muted mb-2 line-clamp-2 text-sm">
          {description || "No description."}
        </p>
      )}

      {meta && (
        <div className="app-text-muted mb-3 flex items-center gap-1 text-xs">
          {meta}
        </div>
      )}

      <button
        onClick={onViewTasks}
        className="text-xs text-(--color-primary) hover:underline"
      >
        {viewTasksLabel}
      </button>
    </div>
  );
}