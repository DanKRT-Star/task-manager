import type { TaskStatus } from "../../types/task";

interface TaskStatusSelectProps {
  status: TaskStatus;
  onChange: (status: TaskStatus) => void;
  className?: string;
}

const statusLabel: Record<TaskStatus, string> = {
  pending: "Pending",
  in_progress: "In Progress",
  done: "Done",
};

export default function TaskStatusSelect({
  status,
  onChange,
  className = "",
}: TaskStatusSelectProps) {
  return (
    <select
      value={status}
      onChange={(e) => onChange(e.target.value as TaskStatus)}
      className={`rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) px-2 py-1 text-xs text-(--color-text) outline-none ${className}`}
    >
      {(Object.keys(statusLabel) as TaskStatus[]).map((s) => (
        <option key={s} value={s}>
          {statusLabel[s]}
        </option>
      ))}
    </select>
  );
}