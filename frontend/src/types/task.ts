import type { Label } from "./label";
import type { Project } from "./project";

export type TaskStatus = "pending" | "in_progress" | "done";

export interface Task {
  taskId: number;
  title: string;
  description: string;
  status: TaskStatus;
  projectId?: number;
  project?: Project;
  epicId?: number;
  milestoneId?: number;
  sprintId?: number;
  userId: number;
  assigneeId?: number;
  labels?: Label[];
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
  milestoneId?: number;
  sprintId?: number;
  assigneeId?: number;
}

export interface UpdateTaskPayload {
  title?: string;
  description?: string;
  status?: TaskStatus;
  deadline?: string;
  epicId?: number;
  milestoneId?: number;
  sprintId?: number;
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