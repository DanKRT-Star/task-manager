export interface Project {
  projectId: number;
  name: string;
  description: string;
  deadline?: string;
  ownerId: number;
  createdAt: string;
  updatedAt: string;
}

export interface ProjectMember {
  projectMemberId: number;
  projectId: number;
  userId: number;
  role: "owner" | "member";
  user?: {
    userId: number;
    userName: string;
    email: string;
  };
  createdAt: string;
}

export interface CreateProjectPayload {
  name: string;
  description?: string;
  deadline?: string;
}

export interface UpdateProjectPayload {
  name?: string;
  description?: string;
  deadline?: string;
}

export interface AddMemberPayload {
  email: string;
}