import type { Task, TaskStatus } from "../../types/task";
import TaskStatusSelect from "./task_status_select";

interface TaskListViewProps {
  tasks: Task[];
  epicTitleById: Map<number, string>;
  milestoneTitleById: Map<number, string>;
  sprintNameById: Map<number, string>;
  onStatusChange: (taskId: number, status: TaskStatus) => void;
  onEdit: (task: Task) => void;
  onDelete: (taskId: number) => void;
  onViewDetail: (task: Task) => void;
}

export default function TaskListView({
  tasks,
  epicTitleById,
  milestoneTitleById,
  sprintNameById,
  onStatusChange,
  onEdit,
  onDelete,
  onViewDetail,
}: TaskListViewProps) {
  if (tasks.length === 0) {
    return (
      <div className="app-card py-12 text-center">
        <p className="app-text-muted">No tasks found.</p>
      </div>
    );
  }

  return (
    <div className="app-card overflow-x-auto">
      <table className="w-full text-left text-sm">
        <thead>
          <tr className="border-b border-(--color-border) text-xs app-text-muted">
            <th className="px-4 py-3 font-medium">Title</th>
            <th className="px-4 py-3 font-medium">Epic</th>
            <th className="px-4 py-3 font-medium">Milestone</th>
            <th className="px-4 py-3 font-medium">Sprint</th>
            <th className="px-4 py-3 font-medium">Labels</th>
            <th className="px-4 py-3 font-medium">Status</th>
            <th className="px-4 py-3 font-medium">Deadline</th>
            <th className="px-4 py-3 font-medium"></th>
          </tr>
        </thead>
        <tbody>
          {tasks.map((task) => (
            <tr
              key={task.taskId}
              className="border-b border-(--color-border) last:border-0"
            >
              <td className="px-4 py-3 font-medium text-(--color-text)">
                {task.title}
              </td>
              <td className="app-text-muted px-4 py-3">
                {task.epicId ? epicTitleById.get(task.epicId) || "-" : "-"}
              </td>
              <td className="app-text-muted px-4 py-3">
                {task.milestoneId ? milestoneTitleById.get(task.milestoneId) || "-" : "-"}
              </td>
              <td className="app-text-muted px-4 py-3">
                {task.sprintId ? sprintNameById.get(task.sprintId) || "-" : "-"}
              </td>
              <td className="px-4 py-3">
                <div className="flex flex-wrap gap-1">
                  {task.labels?.map((label) => (
                    <span
                      key={label.labelId}
                      className="rounded-full px-2 py-0.5 text-[10px] font-medium text-white"
                      style={{ backgroundColor: label.color }}
                    >
                      {label.name}
                    </span>
                  ))}
                </div>
              </td>
              <td className="px-4 py-3">
                <TaskStatusSelect
                  status={task.status}
                  onChange={(status) => onStatusChange(task.taskId, status)}
                />
              </td>
              <td className="app-text-muted px-4 py-3">
                {task.deadline ? task.deadline.slice(0, 10) : "-"}
              </td>
              <td className="px-4 py-3">
                <div className="flex justify-end gap-2">
                  <button
                    onClick={() => onViewDetail(task)}
                    className="text-(--color-muted) hover:text-(--color-text)"
                    title="Comments & activity"
                  >
                    <i className="bx bx-message-square-detail text-base"></i>
                  </button>
                  <button
                    onClick={() => onEdit(task)}
                    className="text-blue-500 hover:text-blue-700"
                  >
                    <i className="bx bx-edit text-base"></i>
                  </button>
                  <button
                    onClick={() => onDelete(task.taskId)}
                    className="text-red-500 hover:text-red-700"
                  >
                    <i className="bx bx-trash text-base"></i>
                  </button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}