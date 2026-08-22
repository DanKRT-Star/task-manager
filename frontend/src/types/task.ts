import type { Label } from "./label";

export type TaskStatus = "pending" | "in_progress" | "done";

export interface Task {
  taskId: number;
  title: string;
  description: string;
  status: TaskStatus;
  projectId?: number;
  epicId?: number;
  milestoneId?: number;
  sprintId?: number;
  labels?: Label[];
  userId: number;
  assigneeId?: number;
  deadline: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateTaskPayload {
  title: string;
  description?: string;
  status?: TaskStatus;
  deadline?: string;
  projectId?: number;
  epicId?: number;
  sprintId?: number;
  milestoneId?: number;
  assigneeId?: number;
}

export interface UpdateTaskPayload {
  title?: string;
  description?: string;
  status?: TaskStatus;
  deadline?: string;
  epicId?: number;
  sprintId?: number;
  milestoneId?: number;
  assigneeId?: number;
}

export interface TaskListResponse {
  data: Task[];
  total: number;
  page: number;
  limit: number;
}

export interface GetTasksParams {
  status?: TaskStatus;
  sort?: "deadline_asc" | "deadline_desc";
  page?: number;
  limit?: number;
}