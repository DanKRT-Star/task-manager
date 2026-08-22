import { useNavigate } from "react-router-dom";
import type { Project } from "../../types/project";

interface ProjectCardProps {
  project: Project;
  onEdit?: (project: Project) => void;
  onDelete?: (projectId: number) => void;
}

export default function ProjectCard({ project, onEdit, onDelete }: ProjectCardProps) {
  const navigate = useNavigate();

  const displayDeadline = project.deadline
    ? project.deadline.slice(0, 10)
    : "No deadline";

  return (
    <div
      onClick={() => navigate(`/projects/${project.projectId}`)}
      className="app-card group relative flex h-full cursor-pointer flex-col p-5 transition-shadow hover:shadow-lg"
    >
      {(onEdit || onDelete) && (
        <div className="absolute right-4 top-4 flex gap-1 rounded border border-(--color-border) bg-(--color-surface) p-1 opacity-0 shadow-sm transition-opacity group-hover:opacity-100">
          {onEdit && (
            <button
              onClick={(e) => {
                e.stopPropagation();
                onEdit(project);
              }}
              className="p-1 text-blue-500 hover:text-blue-700"
              title="Edit"
            >
              <i className="bx bx-edit text-base"></i>
            </button>
          )}
          {onDelete && (
            <button
              onClick={(e) => {
                e.stopPropagation();
                onDelete(project.projectId);
              }}
              className="p-1 text-red-500 hover:text-red-700"
              title="Delete"
            >
              <i className="bx bx-trash text-base"></i>
            </button>
          )}
        </div>
      )}

      <h3
        className="mb-3 line-clamp-2 pr-8 text-lg font-bold text-(--color-text)"
        title={project.name}
      >
        {project.name}
      </h3>

      <p className="app-text-muted mb-5 line-clamp-3 flex-1 text-sm">
        {project.description || "No description provided."}
      </p>

      {/* Footer pinned to the bottom */}
      <div className="mt-auto border-t border-(--color-border) pt-3">
        <div className="flex items-center justify-between">
          <span className="app-text-muted flex items-center text-xs font-medium">
            <i className="bx bx-time-five mr-1.5 text-base"></i>
            {displayDeadline}
          </span>
          <span className="app-text-muted text-xs">
            ID: #{project.projectId}
          </span>
        </div>
      </div>
    </div>
  );
}