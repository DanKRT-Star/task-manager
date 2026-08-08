import { useState, useCallback } from "react";
import { taskService } from "../services/task_service";
import { toast } from "sonner";
import type {
  Task,
  CreateTaskPayload,
  UpdateTaskPayload,
  GetTasksParams,
} from "../types/task";

export function useTask() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchTasks = useCallback(async (params?: GetTasksParams) => {
    setLoading(true);
    setError(null);
    try {
      const res = await taskService.getTasks(params);
      setTasks(res.data);
      setTotal(res.total);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch tasks");
    } finally {
      setLoading(false);
    }
  }, []);

  const createTask = async (payload: CreateTaskPayload) => {
  try {
    const newTask = await taskService.createTask(payload);
    setTasks((prev) => [...prev, newTask]);
    toast.success("Tạo công việc thành công");
    return newTask;
  } catch (err) {
    const message =
      (err as { response?: { data?: { error?: string } } })?.response?.data
        ?.error || "Tạo công việc thất bại";
    toast.error(message);
    throw err;
  }
  };

  const updateTask = async (taskId: number, payload: UpdateTaskPayload) => {
    try {
      const updated = await taskService.updateTask(taskId, payload);
      setTasks((prev) => prev.map((t) => (t.taskId === taskId ? updated : t)));
      toast.success("Cập nhật công việc thành công");
      return updated;
    } catch (err) {
      const message =
        (err as { response?: { data?: { error?: string } } })?.response?.data
          ?.error || "Cập nhật công việc thất bại";
      toast.error(message);
      throw err;
    }
  };

  const deleteTask = async (taskId: number) => {
    try {
      await taskService.deleteTask(taskId);
      setTasks((prev) => prev.filter((t) => t.taskId !== taskId));
      toast.success("Xóa công việc thành công");
    } catch (err) {
      const message =
        (err as { response?: { data?: { error?: string } } })?.response?.data
          ?.error || "Xóa công việc thất bại";
      toast.error(message);
      throw err;
    }
  };

  return {
    tasks,
    total,
    loading,
    error,
    fetchTasks,
    createTask,
    updateTask,
    deleteTask,
  };
}