import api from "./axios";
import type { Label, CreateLabelPayload } from "../types/label";

export const labelService = {
  getProjectLabels: async (projectId: number): Promise<Label[]> => {
    const res = await api.get<{ data: Label[] }>(`/projects/${projectId}/labels`);
    return res.data.data;
  },

  createLabel: async (
    projectId: number,
    payload: CreateLabelPayload
  ): Promise<Label> => {
    const res = await api.post<Label>(`/projects/${projectId}/labels`, payload);
    return res.data;
  },

  deleteLabel: async (projectId: number, labelId: number): Promise<void> => {
    await api.delete(`/projects/${projectId}/labels/${labelId}`);
  },

  getTaskLabels: async (taskId: number): Promise<Label[]> => {
    const res = await api.get<{ data: Label[] }>(`/tasks/${taskId}/labels`);
    return res.data.data;
  },

  attachLabel: async (taskId: number, labelId: number): Promise<void> => {
    await api.post(`/tasks/${taskId}/labels/${labelId}`);
  },

  detachLabel: async (taskId: number, labelId: number): Promise<void> => {
    await api.delete(`/tasks/${taskId}/labels/${labelId}`);
  },
};