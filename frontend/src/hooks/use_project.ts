import { useState, useCallback } from "react";
import { toast } from "sonner";
import { projectService } from "../services/project_service";
import type {
  Project,
  ProjectMember,
  CreateProjectPayload,
  UpdateProjectPayload,
} from "../types/project";
import { getErrorMessage } from "../lib/error";

export function useProject() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [currentProject, setCurrentProject] = useState<Project | null>(null);
  const [members, setMembers] = useState<ProjectMember[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchProjects = useCallback(async () => {
    setLoading(true);
    try {
      const data = await projectService.getProjects();
      setProjects(data);
    } catch (err) {
      toast.error(getErrorMessage(err, "Không tải được danh sách project"));
    } finally {
      setLoading(false);
    }
  }, []);

  // Lấy chi tiết 1 project theo id, dùng cho trang Project Detail
  const fetchProject = useCallback(async (projectId: number) => {
    setLoading(true);
    try {
      const data = await projectService.getProject(projectId);
      setCurrentProject(data);
      return data;
    } catch (err) {
      toast.error(getErrorMessage(err, "Không tải được thông tin project"));
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const createProject = async (payload: CreateProjectPayload) => {
    try {
      const project = await projectService.createProject(payload);
      setProjects((prev) => [...prev, project]);
      toast.success("Tạo project thành công");
      return project;
    } catch (err) {
      toast.error(getErrorMessage(err, "Tạo project thất bại"));
      throw err;
    }
  };

  const updateProject = async (
    projectId: number,
    payload: UpdateProjectPayload
  ) => {
    try {
      const updated = await projectService.updateProject(projectId, payload);
      setProjects((prev) =>
        prev.map((p) => (p.projectId === projectId ? updated : p))
      );
      setCurrentProject((prev) =>
        prev && prev.projectId === projectId ? updated : prev
      );
      toast.success("Cập nhật project thành công");
      return updated;
    } catch (err) {
      toast.error(getErrorMessage(err, "Cập nhật project thất bại"));
      throw err;
    }
  };

  const deleteProject = async (projectId: number) => {
    try {
      await projectService.deleteProject(projectId);
      setProjects((prev) => prev.filter((p) => p.projectId !== projectId));
      toast.success("Xóa project thành công");
    } catch (err) {
      toast.error(getErrorMessage(err, "Xóa project thất bại"));
      throw err;
    }
  };

  const fetchMembers = useCallback(async (projectId: number) => {
    try {
      const data = await projectService.getMembers(projectId);
      setMembers(data);
    } catch (err) {
      toast.error(getErrorMessage(err, "Không tải được danh sách thành viên"));
    }
  }, []);

  const addMember = async (projectId: number, email: string) => {
    try {
      await projectService.addMember(projectId, { email });
      toast.success("Thêm thành viên thành công");
      await fetchMembers(projectId);
    } catch (err) {
      toast.error(getErrorMessage(err, "Thêm thành viên thất bại"));
      throw err;
    }
  };

  const removeMember = async (projectId: number, userId: number) => {
    try {
      await projectService.removeMember(projectId, userId);
      setMembers((prev) => prev.filter((m) => m.userId !== userId));
      toast.success("Xóa thành viên thành công");
    } catch (err) {
      toast.error(getErrorMessage(err, "Xóa thành viên thất bại"));
      throw err;
    }
  };

  return {
    projects,
    currentProject,
    members,
    loading,
    fetchProjects,
    fetchProject,
    createProject,
    updateProject,
    deleteProject,
    fetchMembers,
    addMember,
    removeMember,
  };
}