import api from "./axios";
import type { ActivityLog } from "../types/activity_log";

export const activityLogService = {
  getTaskActivity: async (taskId: number): Promise<ActivityLog[]> => {
    const res = await api.get<{ data: ActivityLog[] }>(`/tasks/${taskId}/activity`);
    return res.data.data;
  },
};