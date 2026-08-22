import { useForm } from "react-hook-form";
import { useEffect } from "react";
import SideFormPanel from "../common/side_form_panel";
import type { Milestone, CreateMilestonePayload } from "../../types/milestone";

interface MilestoneFormPanelProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: CreateMilestonePayload) => Promise<void>;
  editingMilestone?: Milestone | null;
}

export default function MilestoneFormPanel({
  open,
  onClose,
  onSubmit,
  editingMilestone,
}: MilestoneFormPanelProps) {
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<CreateMilestonePayload>({
    defaultValues: { title: "", description: "", dueDate: "" },
  });

  useEffect(() => {
    if (editingMilestone) {
      reset({
        title: editingMilestone.title,
        description: editingMilestone.description || "",
        dueDate: editingMilestone.dueDate?.slice(0, 10) || "",
      });
    } else {
      reset({ title: "", description: "", dueDate: "" });
    }
  }, [editingMilestone, open, reset]);

  const submitHandler = async (data: CreateMilestonePayload) => {
    await onSubmit(data);
    onClose();
  };

  return (
    <SideFormPanel
      open={open}
      onClose={onClose}
      title={editingMilestone ? "Edit Milestone" : "New Milestone"}
      onSubmit={handleSubmit(submitHandler)}
      submitting={isSubmitting}
      submitLabel="Save Milestone"
    >
      <div>
        <label className="mb-1 block text-sm font-medium text-(--color-text)">
          Milestone Title
        </label>
        <input
          type="text"
          {...register("title", { required: "Vui lòng nhập tiêu đề" })}
          className="w-full rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
          placeholder="e.g. Beta Release"
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
          placeholder="What marks this milestone?"
        />
      </div>

      <div>
        <label className="mb-1 block text-sm font-medium text-(--color-text)">
          Due Date
        </label>
        <input
          type="date"
          {...register("dueDate")}
          className="w-full rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
        />
      </div>
    </SideFormPanel>
  );
}