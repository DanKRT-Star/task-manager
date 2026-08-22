import EntitySummaryCard from "../common/entity_summary_card";
import type { Epic } from "../../types/epic";

interface EpicsTabProps {
  epics: Epic[];
  loading: boolean;
  onCreate: () => void;
  onEdit: (epic: Epic) => void;
  onDelete: (epicId: number) => void;
  onViewTasks: (epicId: number) => void;
}

export default function EpicsTab({
  epics,
  loading,
  onCreate,
  onEdit,
  onDelete,
  onViewTasks,
}: EpicsTabProps) {
  return (
    <div>
      <div className="mb-4 flex justify-end">
        <button
          onClick={onCreate}
          className="app-btn-primary flex items-center gap-2 px-4 py-2 text-sm font-medium"
        >
          <i className="bx bx-plus text-lg"></i>
          New Epic
        </button>
      </div>

      {loading ? (
        <div className="app-card py-12 text-center">
          <p className="app-text-muted">Loading epics...</p>
        </div>
      ) : epics.length === 0 ? (
        <div className="app-card py-12 text-center">
          <p className="app-text-muted">No epics yet.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          {epics.map((epic) => (
            <EntitySummaryCard
              key={epic.epicId}
              title={epic.title}
              description={epic.description}
              onEdit={() => onEdit(epic)}
              onDelete={() => onDelete(epic.epicId)}
              onViewTasks={() => onViewTasks(epic.epicId)}
              viewTasksLabel="View tasks in this epic →"
            />
          ))}
        </div>
      )}
    </div>
  );
}