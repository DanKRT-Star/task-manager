export type SprintStatus = "planned" | "active" | "completed";

export interface Sprint {
  sprintId: number;
  projectId: number;
  name: string;
  status: SprintStatus;
  startDate?: string;
  endDate?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateSprintPayload {
  name: string;
  startDate?: string;
  endDate?: string;
}

export interface UpdateSprintPayload {
  name?: string;
  status?: SprintStatus;
  startDate?: string;
  endDate?: string;
}