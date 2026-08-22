// epic.ts
export interface Epic {
  epicId: number;
  projectId: number;
  title: string;
  description: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateEpicPayload {
  title: string;
  description?: string;
}

export interface UpdateEpicPayload {
  title?: string;
  description?: string;
}