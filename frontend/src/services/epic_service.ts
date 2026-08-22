// epic_service.ts
import api from "./axios";
import type { Epic, CreateEpicPayload, UpdateEpicPayload } from "../types/epic";

export const epicService = {
  getProjectEpics: async (projectId: number): Promise<Epic[]> => {
    const res = await api.get<{ data: Epic[] }>(`/projects/${projectId}/epics`);
    return res.data.data;
  },

  createEpic: async (
    projectId: number,
    payload: CreateEpicPayload
  ): Promise<Epic> => {
    const res = await api.post<Epic>(`/projects/${projectId}/epics`, payload);
    return res.data;
  },

  updateEpic: async (
    epicId: number,
    payload: UpdateEpicPayload
  ): Promise<Epic> => {
    const res = await api.put<Epic>(`/epics/${epicId}`, payload);
    return res.data;
  },

  deleteEpic: async (epicId: number): Promise<void> => {
    await api.delete(`/epics/${epicId}`);
  },
};