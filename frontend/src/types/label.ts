export interface Label {
  labelId: number;
  projectId: number;
  name: string;
  color: string;
  createdAt: string;
}

export interface CreateLabelPayload {
  name: string;
  color?: string;
}