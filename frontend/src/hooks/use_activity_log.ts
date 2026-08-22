import { useState, useCallback } from "react";
import { toast } from "sonner";
import { activityLogService } from "../services/activity_log_service";
import type { ActivityLog } from "../types/activity_log";
import { getErrorMessage } from "../lib/error";

export function useActivityLog() {
  const [logs, setLogs] = useState<ActivityLog[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchTaskActivity = useCallback(async (taskId: number) => {
    setLoading(true);
    try {
      const data = await activityLogService.getTaskActivity(taskId);
      setLogs(data);
    } catch (err) {
      toast.error(getErrorMessage(err, "Không tải được lịch sử hoạt động"));
    } finally {
      setLoading(false);
    }
  }, []);

  return { logs, loading, fetchTaskActivity };
}