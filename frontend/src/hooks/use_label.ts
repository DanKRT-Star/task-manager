import { useState, useCallback } from "react";
import { toast } from "sonner";
import { labelService } from "../services/label_service";
import type { Label, CreateLabelPayload } from "../types/label";
import { getErrorMessage } from "../lib/error";

export function useLabel() {
  const [labels, setLabels] = useState<Label[]>([]);
  const [taskLabels, setTaskLabels] = useState<Label[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchLabels = useCallback(async (projectId: number) => {
    setLoading(true);
    try {
      const data = await labelService.getProjectLabels(projectId);
      setLabels(data);
    } catch (err) {
      toast.error(getErrorMessage(err, "Không tải được danh sách label"));
    } finally {
      setLoading(false);
    }
  }, []);

  const fetchTaskLabels = useCallback(async (taskId: number) => {
    try {
      const data = await labelService.getTaskLabels(taskId);
      setTaskLabels(data);
    } catch (err) {
      toast.error(getErrorMessage(err, "Không tải được label của task"));
    }
  }, []);

  const createLabel = async (projectId: number, payload: CreateLabelPayload) => {
    try {
      const label = await labelService.createLabel(projectId, payload);
      setLabels((prev) => [...prev, label]);
      toast.success("Tạo label thành công");
      return label;
    } catch (err) {
      toast.error(getErrorMessage(err, "Tạo label thất bại"));
      throw err;
    }
  };

  // Cần projectId vì route xoá label lồng trong project (/projects/:id/labels/:labelId),
  // không phải route độc lập như epic/milestone/sprint.
  const deleteLabel = async (projectId: number, labelId: number) => {
    try {
      await labelService.deleteLabel(projectId, labelId);
      setLabels((prev) => prev.filter((l) => l.labelId !== labelId));
      toast.success("Xóa label thành công");
    } catch (err) {
      toast.error(getErrorMessage(err, "Xóa label thất bại"));
      throw err;
    }
  };

  const attachLabel = async (taskId: number, labelId: number) => {
    try {
      await labelService.attachLabel(taskId, labelId);
      await fetchTaskLabels(taskId);
      toast.success("Gắn label thành công");
    } catch (err) {
      toast.error(getErrorMessage(err, "Gắn label thất bại"));
      throw err;
    }
  };

  const detachLabel = async (taskId: number, labelId: number) => {
    try {
      await labelService.detachLabel(taskId, labelId);
      setTaskLabels((prev) => prev.filter((l) => l.labelId !== labelId));
      toast.success("Gỡ label thành công");
    } catch (err) {
      toast.error(getErrorMessage(err, "Gỡ label thất bại"));
      throw err;
    }
  };

  return {
    labels,
    taskLabels,
    loading,
    fetchLabels,
    fetchTaskLabels,
    createLabel,
    deleteLabel,
    attachLabel,
    detachLabel,
  };
}