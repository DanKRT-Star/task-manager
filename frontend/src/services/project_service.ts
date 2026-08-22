import api from "./axios";
import type {
  Project,
  ProjectMember,
  CreateProjectPayload,
  UpdateProjectPayload,
  AddMemberPayload,
} from "../types/project";

export const projectService = {
  getProjects: async (): Promise<Project[]> => {
    const res = await api.get<{ data: Project[] }>("/projects");
    return res.data.data;
  },

  getProject: async (projectId: number): Promise<Project> => {
    const res = await api.get<Project>(`/projects/${projectId}`);
    return res.data;
  },

  createProject: async (payload: CreateProjectPayload): Promise<Project> => {
    const res = await api.post<Project>("/projects", payload);
    return res.data;
  },

  updateProject: async (
    projectId: number,
    payload: UpdateProjectPayload
  ): Promise<Project> => {
    const res = await api.put<Project>(`/projects/${projectId}`, payload);
    return res.data;
  },

  deleteProject: async (projectId: number): Promise<void> => {
    await api.delete(`/projects/${projectId}`);
  },

  getMembers: async (projectId: number): Promise<ProjectMember[]> => {
    const res = await api.get<{ data: ProjectMember[] }>(
      `/projects/${projectId}/members`
    );
    return res.data.data;
  },

  addMember: async (
    projectId: number,
    payload: AddMemberPayload
  ): Promise<ProjectMember> => {
    const res = await api.post<ProjectMember>(
      `/projects/${projectId}/members`,
      payload
    );
    return res.data;
  },

  removeMember: async (projectId: number, userId: number): Promise<void> => {
    await api.delete(`/projects/${projectId}/members/${userId}`);
  },
};