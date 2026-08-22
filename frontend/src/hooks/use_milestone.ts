import { useState, useCallback } from "react";
import { toast } from "sonner";
import { milestoneService } from "../services/milestone_service";
import type {
  Milestone,
  CreateMilestonePayload,
  UpdateMilestonePayload,
} from "../types/milestone";
import { getErrorMessage } from "../lib/error";

export function useMilestone() {
  const [milestones, setMilestones] = useState<Milestone[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchMilestones = useCallback(async (projectId: number) => {
    setLoading(true);
    try {
      const data = await milestoneService.getProjectMilestones(projectId);
      setMilestones(data);
    } catch (err) {
      toast.error(getErrorMessage(err, "Không tải được danh sách milestone"));
    } finally {
      setLoading(false);
    }
  }, []);

  const createMilestone = async (
    projectId: number,
    payload: CreateMilestonePayload
  ) => {
    try {
      const milestone = await milestoneService.createMilestone(
        projectId,
        payload
      );
      setMilestones((prev) => [...prev, milestone]);
      toast.success("Tạo milestone thành công");
      return milestone;
    } catch (err) {
      toast.error(getErrorMessage(err, "Tạo milestone thất bại"));
      throw err;
    }
  };

  const updateMilestone = async (
    milestoneId: number,
    payload: UpdateMilestonePayload
  ) => {
    try {
      const updated = await milestoneService.updateMilestone(
        milestoneId,
        payload
      );
      setMilestones((prev) =>
        prev.map((m) => (m.milestoneId === milestoneId ? updated : m))
      );
      toast.success("Cập nhật milestone thành công");
      return updated;
    } catch (err) {
      toast.error(getErrorMessage(err, "Cập nhật milestone thất bại"));
      throw err;
    }
  };

  const deleteMilestone = async (milestoneId: number) => {
    try {
      await milestoneService.deleteMilestone(milestoneId);
      setMilestones((prev) =>
        prev.filter((m) => m.milestoneId !== milestoneId)
      );
      toast.success("Xóa milestone thành công");
    } catch (err) {
      toast.error(getErrorMessage(err, "Xóa milestone thất bại"));
      throw err;
    }
  };

  return {
    milestones,
    loading,
    fetchMilestones,
    createMilestone,
    updateMilestone,
    deleteMilestone,
  };
}