import { useForm } from "react-hook-form";
import { useEffect } from "react";
import SideFormPanel from "../common/side_form_panel";
import type { Task } from "../../types/task";
import type { TaskFormData } from "../../lib/validation";

interface TaskFormPanelProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: TaskFormData) => Promise<void>;
  editingTask?: Task | null;
}

export default function TaskFormPanel({
  open,
  onClose,
  onSubmit,
  editingTask,
}: TaskFormPanelProps) {
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<TaskFormData>({
    defaultValues: { title: "", description: "", status: "pending", deadline: "" },
  });

  useEffect(() => {
    if (editingTask) {
      reset({
        title: editingTask.title,
        description: editingTask.description || "",
        status: editingTask.status,
        deadline: editingTask.deadline?.slice(0, 10) || "",
      });
    } else {
      reset({ title: "", description: "", status: "pending", deadline: "" });
    }
  }, [editingTask, open, reset]);

  const submitHandler = async (data: TaskFormData) => {
    await onSubmit(data);
    onClose();
  };

  return (
    <SideFormPanel
      open={open}
      onClose={onClose}
      title={editingTask ? "Edit Task" : "New Task"}
      onSubmit={handleSubmit(submitHandler)}
      submitting={isSubmitting}
      submitLabel="Save Task"
    >
      <div>
        <label className="mb-1 block text-sm font-medium text-(--color-text)">
          Task Title
        </label>
        <input
          type="text"
          {...register("title", { required: "Vui lòng nhập tiêu đề" })}
          className="w-full rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
          placeholder="e.g. Buy groceries"
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
          rows={4}
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
          Deadline
        </label>
        <input
          type="date"
          {...register("deadline", { required: "Vui lòng chọn deadline" })}
          className="w-full rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
        />
        {errors.deadline && (
          <p className="mt-1 text-xs text-(--color-danger)">
            {errors.deadline.message}
          </p>
        )}
      </div>
    </SideFormPanel>
  );
}