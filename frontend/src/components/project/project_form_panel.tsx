import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect } from "react";
import SideFormPanel from "../common/side_form_panel";
// You will need to add projectSchema and ProjectFormData to your validation file
import { projectSchema, type ProjectFormData } from "../../lib/validation";
import type { Project } from "../../types/project";

interface ProjectFormPanelProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: ProjectFormData) => Promise<void>;
  editingProject?: Project | null;
}

export default function ProjectFormPanel({
  open,
  onClose,
  onSubmit,
  editingProject,
}: ProjectFormPanelProps) {
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<ProjectFormData>({
    resolver: zodResolver(projectSchema),
    defaultValues: { name: "", description: "", deadline: "" },
  });

  useEffect(() => {
    if (editingProject) {
      reset({
        name: editingProject.name,
        description: editingProject.description || "",
        deadline: editingProject.deadline?.slice(0, 10) || "",
      });
    } else {
      reset({ name: "", description: "", deadline: "" });
    }
  }, [editingProject, open, reset]);

  const submitHandler = async (data: ProjectFormData) => {
    await onSubmit(data);
    onClose();
  };

  return (
    <SideFormPanel
      open={open}
      onClose={onClose}
      title={editingProject ? "Edit Project" : "New Project"}
      onSubmit={handleSubmit(submitHandler)}
      submitting={isSubmitting}
      submitLabel="Save Project"
    >
      <div>
        <label className="mb-1 block text-sm font-medium text-(--color-text)">
          Project Name
        </label>
        <input
          type="text"
          {...register("name")}
          className="w-full rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
          placeholder="e.g. Website Redesign"
        />
        {errors.name && (
          <p className="mt-1 text-xs text-(--color-danger)">
            {errors.name.message}
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
          placeholder="What is this project about?"
        />
        {errors.description && (
          <p className="mt-1 text-xs text-(--color-danger)">
            {errors.description.message}
          </p>
        )}
      </div>

      <div>
        <label className="mb-1 block text-sm font-medium text-(--color-text)">
          Deadline
        </label>
        <input
          type="date"
          {...register("deadline")}
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