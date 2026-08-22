import EntitySummaryCard from "../common/entity_summary_card";
import type { Milestone } from "../../types/milestone";

interface MilestonesTabProps {
  milestones: Milestone[];
  loading: boolean;
  onCreate: () => void;
  onEdit: (milestone: Milestone) => void;
  onDelete: (milestoneId: number) => void;
  onViewTasks: (milestoneId: number) => void;
}

export default function MilestonesTab({
  milestones,
  loading,
  onCreate,
  onEdit,
  onDelete,
  onViewTasks,
}: MilestonesTabProps) {
  return (
    <div>
      <div className="mb-4 flex justify-end">
        <button
          onClick={onCreate}
          className="app-btn-primary flex items-center gap-2 px-4 py-2 text-sm font-medium"
        >
          <i className="bx bx-plus text-lg"></i>
          New Milestone
        </button>
      </div>

      {loading ? (
        <div className="app-card py-12 text-center">
          <p className="app-text-muted">Loading milestones...</p>
        </div>
      ) : milestones.length === 0 ? (
        <div className="app-card py-12 text-center">
          <p className="app-text-muted">No milestones yet.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          {milestones.map((m) => (
            <EntitySummaryCard
              key={m.milestoneId}
              title={m.title}
              description={m.description}
              meta={
                <>
                  <i className="bx bx-flag text-base"></i>
                  {m.dueDate ? m.dueDate.slice(0, 10) : "No due date"}
                </>
              }
              onEdit={() => onEdit(m)}
              onDelete={() => onDelete(m.milestoneId)}
              onViewTasks={() => onViewTasks(m.milestoneId)}
              viewTasksLabel="View tasks in this milestone →"
            />
          ))}
        </div>
      )}
    </div>
  );
}