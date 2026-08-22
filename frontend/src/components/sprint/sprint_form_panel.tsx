import { useForm } from "react-hook-form";
import { useEffect } from "react";
import SideFormPanel from "../common/side_form_panel";
import type { Sprint, CreateSprintPayload, SprintStatus } from "../../types/sprint";

interface SprintFormPanelProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: CreateSprintPayload & { status?: SprintStatus }) => Promise<void>;
  editingSprint?: Sprint | null;
}

interface SprintFormValues {
  name: string;
  status: SprintStatus;
  startDate?: string;
  endDate?: string;
}

export default function SprintFormPanel({
  open,
  onClose,
  onSubmit,
  editingSprint,
}: SprintFormPanelProps) {
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<SprintFormValues>({
    defaultValues: { name: "", status: "planned", startDate: "", endDate: "" },
  });

  useEffect(() => {
    if (editingSprint) {
      reset({
        name: editingSprint.name,
        status: editingSprint.status,
        startDate: editingSprint.startDate?.slice(0, 10) || "",
        endDate: editingSprint.endDate?.slice(0, 10) || "",
      });
    } else {
      reset({ name: "", status: "planned", startDate: "", endDate: "" });
    }
  }, [editingSprint, open, reset]);

  const submitHandler = async (data: SprintFormValues) => {
    await onSubmit({
      name: data.name,
      status: editingSprint ? data.status : undefined,
      startDate: data.startDate ? new Date(data.startDate).toISOString() : undefined,
      endDate: data.endDate ? new Date(data.endDate).toISOString() : undefined,
    });
    onClose();
  };

  return (
    <SideFormPanel
      open={open}
      onClose={onClose}
      title={editingSprint ? "Edit Sprint" : "New Sprint"}
      onSubmit={handleSubmit(submitHandler)}
      submitting={isSubmitting}
      submitLabel="Save Sprint"
    >
      <div>
        <label className="mb-1 block text-sm font-medium text-(--color-text)">
          Sprint Name
        </label>
        <input
          type="text"
          {...register("name", { required: "Vui lòng nhập tên sprint" })}
          className="w-full rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
          placeholder="e.g. Sprint 3"
        />
        {errors.name && (
          <p className="mt-1 text-xs text-(--color-danger)">
            {errors.name.message}
          </p>
        )}
      </div>

      {editingSprint && (
        <div>
          <label className="mb-1 block text-sm font-medium text-(--color-text)">
            Status
          </label>
          <select
            {...register("status")}
            className="w-full rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
          >
            <option value="planned">Planned</option>
            <option value="active">Active</option>
            <option value="completed">Completed</option>
          </select>
        </div>
      )}

      <div>
        <label className="mb-1 block text-sm font-medium text-(--color-text)">
          Start Date
        </label>
        <input
          type="date"
          {...register("startDate")}
          className="w-full rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
        />
      </div>

      <div>
        <label className="mb-1 block text-sm font-medium text-(--color-text)">
          End Date
        </label>
        <input
          type="date"
          {...register("endDate")}
          className="w-full rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
        />
      </div>
    </SideFormPanel>
  );
}