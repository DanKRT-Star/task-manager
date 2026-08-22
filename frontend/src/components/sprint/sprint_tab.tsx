import EntitySummaryCard from "../common/entity_summary_card";
import type { Sprint } from "../../types/sprint";

interface SprintsTabProps {
  sprints: Sprint[];
  loading: boolean;
  onCreate: () => void;
  onEdit: (sprint: Sprint) => void;
  onDelete: (sprintId: number) => void;
  onViewTasks: (sprintId: number) => void;
}

export default function SprintsTab({
  sprints,
  loading,
  onCreate,
  onEdit,
  onDelete,
  onViewTasks,
}: SprintsTabProps) {
  return (
    <div>
      <div className="mb-4 flex justify-end">
        <button
          onClick={onCreate}
          className="app-btn-primary flex items-center gap-2 px-4 py-2 text-sm font-medium"
        >
          <i className="bx bx-plus text-lg"></i>
          New Sprint
        </button>
      </div>

      {loading ? (
        <div className="app-card py-12 text-center">
          <p className="app-text-muted">Loading sprints...</p>
        </div>
      ) : sprints.length === 0 ? (
        <div className="app-card py-12 text-center">
          <p className="app-text-muted">No sprints yet.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          {sprints.map((s) => (
            <EntitySummaryCard
              key={s.sprintId}
              title={s.name}
              badge={
                <span className="mb-2 inline-block rounded-full bg-(--color-bg) px-2 py-0.5 text-xs capitalize text-(--color-muted)">
                  {s.status}
                </span>
              }
              meta={
                <>
                  <i className="bx bx-calendar text-base"></i>
                  {s.startDate ? s.startDate.slice(0, 10) : "?"} →{" "}
                  {s.endDate ? s.endDate.slice(0, 10) : "?"}
                </>
              }
              onEdit={() => onEdit(s)}
              onDelete={() => onDelete(s.sprintId)}
              onViewTasks={() => onViewTasks(s.sprintId)}
              viewTasksLabel="View tasks in this sprint →"
            />
          ))}
        </div>
      )}
    </div>
  );
}