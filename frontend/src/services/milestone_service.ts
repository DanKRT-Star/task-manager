// milestone_service.ts
import api from "./axios";
import type {
  Milestone,
  CreateMilestonePayload,
  UpdateMilestonePayload,
} from "../types/milestone";

export const milestoneService = {
  getProjectMilestones: async (projectId: number): Promise<Milestone[]> => {
    const res = await api.get<{ data: Milestone[] }>(
      `/projects/${projectId}/milestones`
    );
    return res.data.data;
  },

  createMilestone: async (
    projectId: number,
    payload: CreateMilestonePayload
  ): Promise<Milestone> => {
    const res = await api.post<Milestone>(
      `/projects/${projectId}/milestones`,
      payload
    );
    return res.data;
  },

  updateMilestone: async (
    milestoneId: number,
    payload: UpdateMilestonePayload
  ): Promise<Milestone> => {
    const res = await api.put<Milestone>(`/milestones/${milestoneId}`, payload);
    return res.data;
  },

  deleteMilestone: async (milestoneId: number): Promise<void> => {
    await api.delete(`/milestones/${milestoneId}`);
  },
};