import type { DraggableProvided } from "@hello-pangea/dnd";
import type { Task, TaskStatus } from "../../types/task";
import TaskStatusSelect from "./task_status_select";

interface TaskCardProps {
  task: Task;
  epicTitle?: string;
  milestoneTitle?: string;
  sprintName?: string;
  onStatusChange: (taskId: number, status: TaskStatus) => void;
  onEdit: (task: Task) => void;
  onDelete: (taskId: number) => void;
  onViewDetail: (task: Task) => void;
  onDropLabel: (taskId: number, labelId: number) => void;
  dragProvided: DraggableProvided;
  isDragging: boolean;
}

export default function TaskCard({
  task,
  epicTitle,
  milestoneTitle,
  sprintName,
  onStatusChange,
  onEdit,
  onDelete,
  onViewDetail,
  onDropLabel,
  dragProvided,
  isDragging,
}: TaskCardProps) {
  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    const labelId = e.dataTransfer.getData("application/x-label-id");
    if (labelId) {
      onDropLabel(task.taskId, Number(labelId));
    }
  };

  const baseStyle = dragProvided.draggableProps.style;
  const cardStyle: React.CSSProperties = isDragging
    ? {
        ...baseStyle,
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        transform: `${(baseStyle as any)?.transform ?? ""} scale(1.03)`,
      }
    : baseStyle ?? {};

  return (
    <div
      // eslint-disable-next-line react-hooks/refs -- required pattern from @hello-pangea/dnd, not a stale-ref access
      ref={dragProvided.innerRef}
      // eslint-disable-next-line react-hooks/refs
      {...dragProvided.draggableProps}
      // eslint-disable-next-line react-hooks/refs
      {...dragProvided.dragHandleProps}
      onDragOver={handleDragOver}
      onDrop={handleDrop}
      style={cardStyle}
      className={`app-card p-4 transition-shadow ${
        isDragging
          ? "shadow-2xl ring-2 ring-(--color-primary)/50"
          : "hover:shadow-md"
      }`}
    >
      <div className="mb-2 flex items-start justify-between gap-3">
        <h4 className="line-clamp-2 text-sm font-semibold text-(--color-text)">
          {task.title}
        </h4>
        <div className="flex shrink-0 gap-2">
          <button
            onClick={() => onViewDetail(task)}
            className="text-(--color-muted) hover:text-(--color-text)"
            title="Comments & activity"
          >
            <i className="bx bx-message-square-detail text-base"></i>
          </button>
          <button
            onClick={() => onEdit(task)}
            className="text-blue-500 hover:text-blue-700"
          >
            <i className="bx bx-edit text-base"></i>
          </button>
          <button
            onClick={() => onDelete(task.taskId)}
            className="text-red-500 hover:text-red-700"
          >
            <i className="bx bx-trash text-base"></i>
          </button>
        </div>
      </div>

      {task.description && (
        <p className="app-text-muted mb-3 line-clamp-2 text-xs">
          {task.description}
        </p>
      )}

      {(epicTitle || milestoneTitle || sprintName) && (
        <div className="mb-2 flex flex-wrap gap-1.5">
          {epicTitle && (
            <span className="rounded-full bg-(--color-bg) px-2 py-0.5 text-xs text-(--color-primary)">
              <i className="bx bx-shape-square mr-1"></i>
              {epicTitle}
            </span>
          )}
          {milestoneTitle && (
            <span className="rounded-full bg-(--color-bg) px-2 py-0.5 text-xs text-(--color-muted)">
              <i className="bx bx-flag mr-1"></i>
              {milestoneTitle}
            </span>
          )}
          {sprintName && (
            <span className="rounded-full bg-(--color-bg) px-2 py-0.5 text-xs text-(--color-muted)">
              <i className="bx bx-run mr-1"></i>
              {sprintName}
            </span>
          )}
        </div>
      )}

      {task.labels && task.labels.length > 0 && (
        <div className="mb-3 flex flex-wrap gap-1">
          {task.labels.map((label) => (
            <span
              key={label.labelId}
              className="rounded-full px-2 py-0.5 text-[10px] font-medium text-white"
              style={{ backgroundColor: label.color }}
            >
              {label.name}
            </span>
          ))}
        </div>
      )}

      <div className="flex items-center justify-between border-t border-(--color-border) pt-2.5">
        <TaskStatusSelect
          status={task.status}
          onChange={(status) => onStatusChange(task.taskId, status)}
        />
        <span className="app-text-muted text-xs">
          {task.deadline ? task.deadline.slice(0, 10) : "No deadline"}
        </span>
      </div>
    </div>
  );
}