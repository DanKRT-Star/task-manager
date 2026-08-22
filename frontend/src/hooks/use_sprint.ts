import { useState, useCallback } from "react";
import { toast } from "sonner";
import { sprintService } from "../services/sprint_service";
import type {
  Sprint,
  CreateSprintPayload,
  UpdateSprintPayload,
} from "../types/sprint";
import { getErrorMessage } from "../lib/error";

export function useSprint() {
  const [sprints, setSprints] = useState<Sprint[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchSprints = useCallback(async (projectId: number) => {
    setLoading(true);
    try {
      const data = await sprintService.getProjectSprints(projectId);
      setSprints(data);
    } catch (err) {
      toast.error(getErrorMessage(err, "Không tải được danh sách sprint"));
    } finally {
      setLoading(false);
    }
  }, []);

  const createSprint = async (
    projectId: number,
    payload: CreateSprintPayload
  ) => {
    try {
      const sprint = await sprintService.createSprint(projectId, payload);
      setSprints((prev) => [...prev, sprint]);
      toast.success("Tạo sprint thành công");
      return sprint;
    } catch (err) {
      toast.error(getErrorMessage(err, "Tạo sprint thất bại"));
      throw err;
    }
  };

  const updateSprint = async (
    sprintId: number,
    payload: UpdateSprintPayload
  ) => {
    try {
      const updated = await sprintService.updateSprint(sprintId, payload);
      setSprints((prev) =>
        prev.map((s) => (s.sprintId === sprintId ? updated : s))
      );
      toast.success("Cập nhật sprint thành công");
      return updated;
    } catch (err) {
      toast.error(getErrorMessage(err, "Cập nhật sprint thất bại"));
      throw err;
    }
  };

  const deleteSprint = async (sprintId: number) => {
    try {
      await sprintService.deleteSprint(sprintId);
      setSprints((prev) => prev.filter((s) => s.sprintId !== sprintId));
      toast.success("Xóa sprint thành công");
    } catch (err) {
      toast.error(getErrorMessage(err, "Xóa sprint thất bại"));
      throw err;
    }
  };

  return {
    sprints,
    loading,
    fetchSprints,
    createSprint,
    updateSprint,
    deleteSprint,
  };
}