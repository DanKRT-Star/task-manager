import { useEffect, useState } from "react";
import ProjectCard from "../components/project/project_card";
import ProjectFormPanel from "../components/project/project_form_panel";
import { useProject } from "../hooks/use_project";
import type { Project } from "../types/project";
import type { ProjectFormData } from "../lib/validation";

export default function ProjectsPage() {
  const { projects, loading, fetchProjects, createProject, updateProject, deleteProject } = useProject();

  const [searchQuery, setSearchQuery] = useState("");
  const [panelOpen, setPanelOpen] = useState(false);
  const [editingProject, setEditingProject] = useState<Project | null>(null);

  useEffect(() => {
    fetchProjects();
  }, [fetchProjects]);

  const handleCreate = () => {
    setEditingProject(null);
    setPanelOpen(true);
  };

  const handleEdit = (project: Project) => {
    setEditingProject(project);
    setPanelOpen(true);
  };

  const handleDelete = async (projectId: number) => {
    if (!confirm("Are you sure you want to delete this project?")) return;
    try {
      await deleteProject(projectId);
    } catch {
      // Toast is already handled inside the hook
    }
  };

  const handleSubmit = async (data: ProjectFormData) => {
    try {
      const payload = {
        ...data,
        deadline: data.deadline ? new Date(data.deadline).toISOString() : undefined,
      };

      if (editingProject) {
        await updateProject(editingProject.projectId, payload);
      } else {
        await createProject(payload);
      }
      fetchProjects();
    } catch {
      // Toast is already handled inside the hook
    }
  };

  const filteredProjects = projects?.filter((p: Project) =>
    p.name.toLowerCase().includes(searchQuery.toLowerCase())
  ) || [];

  return (
    <div className="mx-auto max-w-6xl p-4 sm:p-6">
      <div className="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <h1 className="app-panel-title">Projects</h1>

        <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
          <div className="relative">
            <i className="bx bx-search absolute left-3 top-1/2 -translate-y-1/2 text-(--color-muted)"></i>
            <input
              type="text"
              placeholder="Search projects..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full rounded-(--radius-button) border border-(--color-border) bg-(--color-surface) py-2 pl-9 pr-4 text-(--text-sm) outline-none focus:border-(--color-primary) sm:w-64"
            />
          </div>

          <button
            onClick={handleCreate}
            className="app-btn-primary flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium"
          >
            <i className="bx bx-plus text-lg"></i>
            New Project
          </button>
        </div>
      </div>

      {loading ? (
        <div className="app-card py-12 text-center">
          <p className="app-text-muted">Loading projects...</p>
        </div>
      ) : filteredProjects.length === 0 ? (
        <div className="app-card py-12 text-center">
          <i className="bx bx-folder-open mb-2 text-4xl text-(--color-muted)"></i>
          <p className="app-text-muted">No projects found.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {filteredProjects.map((project: Project) => (
            <ProjectCard
              key={project.projectId}
              project={project}
              onEdit={handleEdit}
              onDelete={handleDelete}
            />
          ))}
        </div>
      )}

      <ProjectFormPanel
        open={panelOpen}
        onClose={() => setPanelOpen(false)}
        onSubmit={handleSubmit}
        editingProject={editingProject}
      />
    </div>
  );
}