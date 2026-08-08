import api from "./axios";
import type {
  Task,
  CreateTaskPayload,
  UpdateTaskPayload,
  TaskListResponse,
  GetTasksParams,
} from "../types/task";

export const taskService = {
  getTasks: async (params?: GetTasksParams): Promise<TaskListResponse> => {
    const res = await api.get<TaskListResponse>("/tasks", { params });
    return res.data;
  },

  createTask: async (payload: CreateTaskPayload): Promise<Task> => {
    const res = await api.post<Task>("/tasks", payload);
    return res.data;
  },

  updateTask: async (
    taskId: number,
    payload: UpdateTaskPayload
  ): Promise<Task> => {
    const res = await api.put<Task>(`/tasks/${taskId}`, payload);
    return res.data;
  },

  deleteTask: async (taskId: number): Promise<{ message: string }> => {
    const res = await api.delete<{ message: string }>(`/tasks/${taskId}`);
    return res.data;
  },
};