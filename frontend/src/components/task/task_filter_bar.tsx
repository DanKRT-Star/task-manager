import { useEffect, useRef, useState } from "react";
import type { Epic } from "../../types/epic";
import type { Milestone } from "../../types/milestone";
import type { Sprint } from "../../types/sprint";
import type { Label } from "../../types/label";
import type { TaskStatus } from "../../types/task";

export type TaskView = "kanban" | "list";

interface TaskFilterBarProps {
  epics: Epic[];
  milestones: Milestone[];
  sprints: Sprint[];
  labels: Label[];
  epicFilter: number | null;
  milestoneFilter: number | null;
  sprintFilter: number | null;
  labelFilter: number[];
  statusFilter: TaskStatus | null;
  view: TaskView;
  onEpicFilterChange: (epicId: number | null) => void;
  onMilestoneFilterChange: (milestoneId: number | null) => void;
  onSprintFilterChange: (sprintId: number | null) => void;
  onLabelFilterToggle: (labelId: number) => void;
  onStatusFilterChange: (status: TaskStatus | null) => void;
  onViewChange: (view: TaskView) => void;
  onNewTask: () => void;
}

export default function TaskFilterBar({
  epics,
  milestones,
  sprints,
  labels,
  epicFilter,
  milestoneFilter,
  sprintFilter,
  labelFilter,
  statusFilter,
  view,
  onEpicFilterChange,
  onMilestoneFilterChange,
  onSprintFilterChange,
  onLabelFilterToggle,
  onStatusFilterChange,
  onViewChange,
  onNewTask,
}: TaskFilterBarProps) {
  const [labelDropdownOpen, setLabelDropdownOpen] = useState(false);
  const labelDropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (
        labelDropdownRef.current &&
        !labelDropdownRef.current.contains(e.target as Node)
      ) {
        setLabelDropdownOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const labelFilterSummary = () => {
    if (labelFilter.length === 0) return "All Labels";
    const first = labels.find((l) => l.labelId === labelFilter[0]);
    const firstName = first?.name || "Label";
    return labelFilter.length === 1 ? firstName : `${firstName} (+${labelFilter.length - 1})`;
  };

  const handleLabelDragStart = (e: React.DragEvent, labelId: number) => {
    e.dataTransfer.setData("application/x-label-id", String(labelId));
    e.dataTransfer.effectAllowed = "copy";
  };

  return (
    <div className="mb-4 flex flex-col gap-3">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex flex-wrap items-center gap-2">
          <select
            value={epicFilter ?? ""}
            onChange={(e) =>
              onEpicFilterChange(e.target.value ? Number(e.target.value) : null)
            }
            className="rounded-(--radius-button) border border-(--color-border) bg-(--color-surface) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
          >
            <option value="">All Epics</option>
            {epics.map((epic) => (
              <option key={epic.epicId} value={epic.epicId}>
                {epic.title}
              </option>
            ))}
          </select>

          <select
            value={milestoneFilter ?? ""}
            onChange={(e) =>
              onMilestoneFilterChange(e.target.value ? Number(e.target.value) : null)
            }
            className="rounded-(--radius-button) border border-(--color-border) bg-(--color-surface) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
          >
            <option value="">All Milestones</option>
            {milestones.map((m) => (
              <option key={m.milestoneId} value={m.milestoneId}>
                {m.title}
              </option>
            ))}
          </select>

          <select
            value={sprintFilter ?? ""}
            onChange={(e) =>
              onSprintFilterChange(e.target.value ? Number(e.target.value) : null)
            }
            className="rounded-(--radius-button) border border-(--color-border) bg-(--color-surface) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
          >
            <option value="">All Sprints</option>
            {sprints.map((s) => (
              <option key={s.sprintId} value={s.sprintId}>
                {s.name}
              </option>
            ))}
          </select>

          <div className="relative" ref={labelDropdownRef}>
            <button
              type="button"
              onClick={() => setLabelDropdownOpen((o) => !o)}
              className="flex items-center gap-1 rounded-(--radius-button) border border-(--color-border) bg-(--color-surface) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
            >
              {labelFilterSummary()}
              <i className="bx bx-chevron-down text-base"></i>
            </button>

            {labelDropdownOpen && (
              <div className="absolute left-0 z-30 mt-1 w-52 rounded-(--radius-button) border border-(--color-border) bg-(--color-surface) p-2 shadow-lg">
                {labels.length === 0 ? (
                  <p className="app-text-muted px-2 py-1.5 text-xs">No labels yet.</p>
                ) : (
                  labels.map((label) => (
                    <label
                      key={label.labelId}
                      className="flex cursor-pointer items-center gap-2 rounded-(--radius-button) px-2 py-1.5 text-sm hover:bg-(--color-bg)"
                    >
                      <input
                        type="checkbox"
                        checked={labelFilter.includes(label.labelId)}
                        onChange={() => onLabelFilterToggle(label.labelId)}
                      />
                      <span
                        className="h-2.5 w-2.5 shrink-0 rounded-full"
                        style={{ backgroundColor: label.color }}
                      ></span>
                      <span className="text-(--color-text)">{label.name}</span>
                    </label>
                  ))
                )}
              </div>
            )}
          </div>

          <select
            value={statusFilter ?? ""}
            onChange={(e) =>
              onStatusFilterChange(e.target.value ? (e.target.value as TaskStatus) : null)
            }
            className="rounded-(--radius-button) border border-(--color-border) bg-(--color-surface) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
          >
            <option value="">All Statuses</option>
            <option value="pending">Pending</option>
            <option value="in_progress">In Progress</option>
            <option value="done">Done</option>
          </select>
        </div>

        <div className="flex items-center gap-2">
          <div className="flex rounded-(--radius-button) border border-(--color-border) p-0.5">
            <button
              onClick={() => onViewChange("kanban")}
              className={`flex items-center gap-1 rounded-(--radius-button) px-3 py-1.5 text-xs font-medium ${
                view === "kanban"
                  ? "bg-(--color-primary) text-white"
                  : "app-text-muted hover:text-(--color-text)"
              }`}
            >
              <i className="bx bx-columns"></i>
              Board
            </button>
            <button
              onClick={() => onViewChange("list")}
              className={`flex items-center gap-1 rounded-(--radius-button) px-3 py-1.5 text-xs font-medium ${
                view === "list"
                  ? "bg-(--color-primary) text-white"
                  : "app-text-muted hover:text-(--color-text)"
              }`}
            >
              <i className="bx bx-list-ul"></i>
              List
            </button>
          </div>

          <button
            onClick={onNewTask}
            className="app-btn-primary flex items-center gap-2 px-4 py-2 text-sm font-medium"
          >
            <i className="bx bx-plus text-lg"></i>
            New Task
          </button>
        </div>
      </div>

      {labels.length > 0 && (
        <div className="flex flex-wrap items-center gap-1.5">
          <span className="app-text-muted mr-1 flex items-center gap-1 text-xs">
            <i className="bx bx-move text-sm"></i>
            Drag onto a task to attach:
          </span>
          {labels.map((label) => (
            <span
              key={label.labelId}
              draggable
              onDragStart={(e) => handleLabelDragStart(e, label.labelId)}
              className="cursor-grab select-none rounded-full px-3 py-1 text-xs font-medium text-white active:cursor-grabbing"
              style={{ backgroundColor: label.color }}
              title="Drag onto a task card"
            >
              {label.name}
            </span>
          ))}
        </div>
      )}
    </div>
  );
}