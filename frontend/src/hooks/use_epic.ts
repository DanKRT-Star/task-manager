// use_epic.ts
import { useState, useCallback } from "react";
import { toast } from "sonner";
import { epicService } from "../services/epic_service";
import type { Epic, CreateEpicPayload, UpdateEpicPayload } from "../types/epic";
import { getErrorMessage } from "../lib/error"

export function useEpic() {
  const [epics, setEpics] = useState<Epic[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchEpics = useCallback(async (projectId: number) => {
    setLoading(true);
    try {
      const data = await epicService.getProjectEpics(projectId);
      setEpics(data);
    } catch (err) {
      toast.error(getErrorMessage(err, "Không tải được danh sách epic"));
    } finally {
      setLoading(false);
    }
  }, []);

  const createEpic = async (projectId: number, payload: CreateEpicPayload) => {
    try {
      const epic = await epicService.createEpic(projectId, payload);
      setEpics((prev) => [...prev, epic]);
      toast.success("Tạo epic thành công");
      return epic;
    } catch (err) {
      toast.error(getErrorMessage(err, "Tạo epic thất bại"));
      throw err;
    }
  };

  const updateEpic = async (epicId: number, payload: UpdateEpicPayload) => {
    try {
      const updated = await epicService.updateEpic(epicId, payload);
      setEpics((prev) => prev.map((e) => (e.epicId === epicId ? updated : e)));
      toast.success("Cập nhật epic thành công");
      return updated;
    } catch (err) {
      toast.error(getErrorMessage(err, "Cập nhật epic thất bại"));
      throw err;
    }
  };

  const deleteEpic = async (epicId: number) => {
    try {
      await epicService.deleteEpic(epicId);
      setEpics((prev) => prev.filter((e) => e.epicId !== epicId));
      toast.success("Xóa epic thành công");
    } catch (err) {
      toast.error(getErrorMessage(err, "Xóa epic thất bại"));
      throw err;
    }
  };

  return { epics, loading, fetchEpics, createEpic, updateEpic, deleteEpic };
}