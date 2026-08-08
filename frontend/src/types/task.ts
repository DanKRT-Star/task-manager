export type TaskStatus = "pending" | "in_progress" | "done";

export interface Task {
  taskId: number;
  title: string;
  description: string;
  status: TaskStatus;
  userId: number;
  deadline: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateTaskPayload {
  title: string;
  description?: string;
  status?: TaskStatus;
  deadline?: string;
}

export interface UpdateTaskPayload {
  title?: string;
  description?: string;
  status?: TaskStatus;
  deadline?: string;
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