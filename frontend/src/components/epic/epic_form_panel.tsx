import { useForm } from "react-hook-form";
import { useEffect } from "react";
import SideFormPanel from "../common/side_form_panel";
import type { Epic, CreateEpicPayload } from "../../types/epic";

interface EpicFormPanelProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: CreateEpicPayload) => Promise<void>;
  editingEpic?: Epic | null;
}

export default function EpicFormPanel({
  open,
  onClose,
  onSubmit,
  editingEpic,
}: EpicFormPanelProps) {
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<CreateEpicPayload>({
    defaultValues: { title: "", description: "" },
  });

  useEffect(() => {
    if (editingEpic) {
      reset({
        title: editingEpic.title,
        description: editingEpic.description || "",
      });
    } else {
      reset({ title: "", description: "" });
    }
  }, [editingEpic, open, reset]);

  const submitHandler = async (data: CreateEpicPayload) => {
    await onSubmit(data);
    onClose();
  };

  return (
    <SideFormPanel
      open={open}
      onClose={onClose}
      title={editingEpic ? "Edit Epic" : "New Epic"}
      onSubmit={handleSubmit(submitHandler)}
      submitting={isSubmitting}
      submitLabel="Save Epic"
    >
      <div>
        <label className="mb-1 block text-sm font-medium text-(--color-text)">
          Epic Title
        </label>
        <input
          type="text"
          {...register("title", { required: "Vui lòng nhập tiêu đề" })}
          className="w-full rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
          placeholder="e.g. User Authentication"
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
          placeholder="What does this epic cover?"
        />
      </div>
    </SideFormPanel>
  );
}