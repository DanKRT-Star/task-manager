import { useForm } from "react-hook-form";
import { useEffect, useState } from "react";
import type { Task, CreateTaskPayload, TaskStatus } from "../../types/task";
import type { Epic } from "../../types/epic";
import type { Milestone } from "../../types/milestone";
import type { Sprint } from "../../types/sprint";
import type { Label } from "../../types/label";
import type { ProjectMember } from "../../types/project";

interface ProjectTaskFormPanelProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: CreateTaskPayload, selectedLabelIds: number[]) => Promise<void>;
  editingTask?: Task | null;
  epics: Epic[];
  milestones: Milestone[];
  sprints: Sprint[];
  labels: Label[];
  members: ProjectMember[];
  defaultEpicId?: number | null;
  defaultMilestoneId?: number | null;
  defaultSprintId?: number | null;
}

interface TaskFormValues {
  title: string;
  description?: string;
  status: TaskStatus;
  deadline?: string;
  epicId?: string;
  milestoneId?: string;
  sprintId?: string;
  assigneeId?: string;
}

export default function ProjectTaskFormPanel({
  open,
  onClose,
  onSubmit,
  editingTask,
  epics,
  milestones,
  sprints,
  labels,
  members,
  defaultEpicId,
  defaultMilestoneId,
  defaultSprintId,
}: ProjectTaskFormPanelProps) {
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<TaskFormValues>({
    defaultValues: { title: "", description: "", status: "pending" },
  });

  const [selectedLabelIds, setSelectedLabelIds] = useState<number[]>([]);

  useEffect(() => {
    if (editingTask) {
      reset({
        title: editingTask.title,
        description: editingTask.description || "",
        status: editingTask.status,
        deadline: editingTask.deadline?.slice(0, 10) || "",
        epicId: editingTask.epicId ? String(editingTask.epicId) : "",
        milestoneId: editingTask.milestoneId ? String(editingTask.milestoneId) : "",
        sprintId: editingTask.sprintId ? String(editingTask.sprintId) : "",
        assigneeId: editingTask.assigneeId ? String(editingTask.assigneeId) : "",
      });
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setSelectedLabelIds(editingTask.labels?.map((l) => l.labelId) || []);
    } else {
      reset({
        title: "",
        description: "",
        status: "pending",
        deadline: "",
        epicId: defaultEpicId ? String(defaultEpicId) : "",
        milestoneId: defaultMilestoneId ? String(defaultMilestoneId) : "",
        sprintId: defaultSprintId ? String(defaultSprintId) : "",
        assigneeId: "",
      });
      setSelectedLabelIds([]);
    }
  }, [editingTask, open, defaultEpicId, defaultMilestoneId, defaultSprintId, reset]);

  const toggleLabel = (labelId: number) => {
    setSelectedLabelIds((prev) =>
      prev.includes(labelId) ? prev.filter((id) => id !== labelId) : [...prev, labelId]
    );
  };

  const submitHandler = async (data: TaskFormValues) => {
    await onSubmit(
      {
        title: data.title,
        description: data.description,
        status: data.status,
        deadline: data.deadline ? new Date(data.deadline).toISOString() : undefined,
        epicId: data.epicId ? Number(data.epicId) : undefined,
        milestoneId: data.milestoneId ? Number(data.milestoneId) : undefined,
        sprintId: data.sprintId ? Number(data.sprintId) : undefined,
        assigneeId: data.assigneeId ? Number(data.assigneeId) : undefined,
      },
      selectedLabelIds
    );
    onClose();
  };

  if (!open) return null;

  return (
    <>
      <div
        className="fixed inset-0 z-10 bg-black/20 transition-opacity"
        onClick={onClose}
      />

      <div className="fixed right-0 top-0 z-20 flex h-full w-80 flex-col overflow-y-auto border-l border-(--color-border) bg-(--color-surface) p-6 shadow-lg transition-transform">
        <div className="mb-6 flex items-center justify-between">
          <h2 className="app-panel-title text-lg">
            {editingTask ? "Edit Task" : "New Task"}
          </h2>
          <button
            onClick={onClose}
            className="text-(--color-muted) hover:text-(--color-text)"
          >
            <i className="bx bx-x text-2xl"></i>
          </button>
        </div>

        <form
          onSubmit={handleSubmit(submitHandler)}
          className="flex h-full flex-col space-y-4"
        >
          <div>
            <label className="mb-1 block text-sm font-medium text-(--color-text)">
              Title
            </label>
            <input
              type="text"
              {...register("title", { required: "Vui lòng nhập tiêu đề" })}
              className="w-full rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
              placeholder="e.g. Implement login API"
            />
            {errors.title && (
              <p className="mt-1 text-xs text-(--color-danger)">
                {errors.title.message}
              </p>
            )}
          </div>

          <div>
            <label className="mb-1 block text-sm font-medium text-(--color-text)">
              Description
            </label>
            <textarea
              {...register("description")}
              rows={3}
              className="w-full resize-none rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
            />
          </div>

          <div>
            <label className="mb-1 block text-sm font-medium text-(--color-text)">
              Status
            </label>
            <select
              {...register("status")}
              className="w-full rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
            >
              <option value="pending">Pending</option>
              <option value="in_progress">In Progress</option>
              <option value="done">Done</option>
            </select>
          </div>

          <div>
            <label className="mb-1 block text-sm font-medium text-(--color-text)">
              Epic
            </label>
            <select
              {...register("epicId")}
              className="w-full rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
            >
              <option value="">None</option>
              {epics.map((epic) => (
                <option key={epic.epicId} value={epic.epicId}>
                  {epic.title}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="mb-1 block text-sm font-medium text-(--color-text)">
              Milestone
            </label>
            <select
              {...register("milestoneId")}
              className="w-full rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
            >
              <option value="">None</option>
              {milestones.map((m) => (
                <option key={m.milestoneId} value={m.milestoneId}>
                  {m.title}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="mb-1 block text-sm font-medium text-(--color-text)">
              Sprint
            </label>
            <select
              {...register("sprintId")}
              className="w-full rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
            >
              <option value="">None</option>
              {sprints.map((s) => (
                <option key={s.sprintId} value={s.sprintId}>
                  {s.name}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="mb-1 block text-sm font-medium text-(--color-text)">
              Assignee
            </label>
            <select
              {...register("assigneeId")}
              className="w-full rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
            >
              <option value="">Unassigned</option>
              {members.map((m) => (
                <option key={m.userId} value={m.userId}>
                  {m.user?.userName || `User #${m.userId}`}
                </option>
              ))}
            </select>
          </div>

          {labels.length > 0 && (
            <div>
              <label className="mb-1 block text-sm font-medium text-(--color-text)">
                Labels
              </label>
              <div className="flex flex-wrap gap-1.5">
                {labels.map((label) => {
                  const active = selectedLabelIds.includes(label.labelId);
                  return (
                    <button
                      key={label.labelId}
                      type="button"
                      onClick={() => toggleLabel(label.labelId)}
                      className="rounded-full px-2.5 py-1 text-xs font-medium text-white transition-opacity"
                      style={{ backgroundColor: label.color, opacity: active ? 1 : 0.35 }}
                    >
                      {label.name}
                    </button>
                  );
                })}
              </div>
            </div>
          )}

          <div>
            <label className="mb-1 block text-sm font-medium text-(--color-text)">
              Deadline
            </label>
            <input
              type="date"
              {...register("deadline")}
              className="w-full rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
            />
          </div>

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
              disabled={isSubmitting}
              className="app-btn-primary px-4 py-2 text-sm font-medium disabled:opacity-50"
            >
              {isSubmitting ? "Saving..." : "Save Task"}
            </button>
          </div>
        </form>
      </div>
    </>
  );
}