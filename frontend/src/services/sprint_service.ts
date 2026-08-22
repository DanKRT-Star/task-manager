import api from "./axios";
import type {
  Sprint,
  CreateSprintPayload,
  UpdateSprintPayload,
} from "../types/sprint";

export const sprintService = {
  getProjectSprints: async (projectId: number): Promise<Sprint[]> => {
    const res = await api.get<{ data: Sprint[] }>(`/projects/${projectId}/sprints`);
    return res.data.data;
  },

  createSprint: async (
    projectId: number,
    payload: CreateSprintPayload
  ): Promise<Sprint> => {
    const res = await api.post<Sprint>(`/projects/${projectId}/sprints`, payload);
    return res.data;
  },

  updateSprint: async (
    sprintId: number,
    payload: UpdateSprintPayload
  ): Promise<Sprint> => {
    const res = await api.put<Sprint>(`/sprints/${sprintId}`, payload);
    return res.data;
  },

  deleteSprint: async (sprintId: number): Promise<void> => {
    await api.delete(`/sprints/${sprintId}`);
  },
};