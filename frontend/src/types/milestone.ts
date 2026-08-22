export interface Milestone {
  milestoneId: number;
  projectId: number;
  title: string;
  description: string;
  dueDate?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateMilestonePayload {
  title: string;
  description?: string;
  dueDate?: string;
}

export interface UpdateMilestonePayload {
  title?: string;
  description?: string;
  dueDate?: string;
}
