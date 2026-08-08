import { useEffect, useState } from "react";
import { useAuth } from "../hooks/use_auth";
import { useTask } from "../hooks/use_task";
import StatusBadge from "../components/status_badge";
import TaskFormPanel from "../components/task_form_panel";
import EmptyTasksRow from "../components/empty_tasks_row";
import AppHeader from "../components/app_header";
import type { Task, TaskStatus } from "../types/task";
import type { TaskFormData } from "../lib/validation";

const PAGE_SIZE = 10;

export default function TasksPage() {
  const { user, logout } = useAuth();
  const { tasks, total, loading, fetchTasks, createTask, updateTask, deleteTask } =
    useTask();

  const [page, setPage] = useState(1);
  const [statusFilter, setStatusFilter] = useState<TaskStatus | "">("");
  const [sort, setSort] = useState<"deadline_asc" | "deadline_desc">(
    "deadline_asc"
  );
  const [panelOpen, setPanelOpen] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);

  const totalPages = Math.max(1, Math.ceil(total / PAGE_SIZE));

  useEffect(() => {
    fetchTasks({
      page,
      limit: PAGE_SIZE,
      status: statusFilter || undefined,
      sort,
    });
  }, [page, statusFilter, sort, fetchTasks]);

  const handleCreate = () => {
    setEditingTask(null);
    setPanelOpen(true);
  };

  const handleEdit = (task: Task) => {
    setEditingTask(task);
    setPanelOpen(true);
  };

  const handleDelete = async (taskId: number) => {
    if (!confirm("Bạn chắc chắn muốn xóa công việc này?")) return;
    try {
      await deleteTask(taskId);
    } catch {
      // toast đã xử lý trong hook
    }
  };

  const handleSubmit = async (data: TaskFormData) => {
    try {
      const payload = {
        ...data,
        deadline: new Date(data.deadline).toISOString(),
      };
      if (editingTask) {
        await updateTask(editingTask.taskId, payload);
      } else {
        await createTask(payload);
      }
      fetchTasks({ page, limit: PAGE_SIZE, status: statusFilter || undefined, sort });
    } catch {
      // toast đã xử lý trong hook
    }
  };

  return (
    <div className="app-shell">
      <AppHeader user={user} onLogout={logout} />

      <main className="relative mx-auto max-w-6xl p-4 sm:p-6">
        <div className="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <h1 className="app-panel-title">My Tasks</h1>

          <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
            <select
              value={sort}
              onChange={(e) =>
                setSort(e.target.value as "deadline_asc" | "deadline_desc")
              }
              className="app-input-select"
            >
              <option value="deadline_asc">Sort: Deadline ↑</option>
              <option value="deadline_desc">Sort: Deadline ↓</option>
            </select>

            <select
              value={statusFilter}
              onChange={(e) => {
                setStatusFilter(e.target.value as TaskStatus | "");
                setPage(1);
              }}
              className="app-input-select"
            >
              <option value="">Filter: All</option>
              <option value="pending">Pending</option>
              <option value="in_progress">In Progress</option>
              <option value="done">Done</option>
            </select>

            <button
              onClick={handleCreate}
              className="app-btn-primary px-4 py-2 text-sm font-medium"
            >
              + Add Task
            </button>
          </div>
        </div>

        <div className="app-card relative overflow-hidden">
          {loading ? (
            <div className="px-4 py-8 text-center text-(--color-muted)">
              Đang tải...
            </div>
          ) : tasks.length === 0 ? (
            <EmptyTasksRow />
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full min-w-160 text-left text-sm">
                <thead className="border-b border-(--color-border) bg-gray-50 text-(--color-muted)">
                  <tr>
                    <th className="px-4 py-3 font-medium">Status</th>
                    <th className="px-4 py-3 font-medium">Task Title</th>
                    <th className="px-4 py-3 font-medium">Description</th>
                    <th className="px-4 py-3 font-medium">Deadline</th>
                    <th className="px-4 py-3 font-medium">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {tasks.map((task) => (
                    <tr key={task.taskId} className="border-b border-gray-100">
                      <td className="px-4 py-3">
                        <StatusBadge status={task.status} />
                      </td>
                      <td className="px-4 py-3 font-medium text-(--color-text)">
                        {task.title}
                      </td>
                      <td className="max-w-xs truncate px-4 py-3 text-(--color-muted)">
                        {task.description || "—"}
                      </td>
                      <td className="px-4 py-3 text-(--color-muted)">
                        {task.deadline?.slice(0, 10)}
                      </td>
                      <td className="px-4 py-3">
                        <div className="flex gap-2">
                          <button
                            onClick={() => handleEdit(task)}
                            className="text-blue-500 hover:text-blue-700"
                            title="Sửa"
                            aria-label="Sửa task"
                          >
                            <i className="bx bx-edit text-xl" />
                          </button>
                          <button
                            onClick={() => handleDelete(task.taskId)}
                            className="text-red-500 hover:text-red-700"
                            title="Xóa"
                            aria-label="Xóa task"
                          >
                            <i className="bx bx-trash text-xl" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}

          <TaskFormPanel
            open={panelOpen}
            onClose={() => setPanelOpen(false)}
            onSubmit={handleSubmit}
            editingTask={editingTask}
          />
        </div>

        <div className="mt-4 flex flex-col gap-2 text-sm text-(--color-muted) sm:flex-row sm:items-center sm:justify-between">
          <span>
            Showing {tasks.length === 0 ? 0 : (page - 1) * PAGE_SIZE + 1}-
            {Math.min(page * PAGE_SIZE, total)} of {total} tasks
          </span>
          <div className="flex items-center gap-1">
            <button
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page === 1}
              className="rounded border border-(--color-border) px-2 py-1 disabled:opacity-40"
            >
              ‹
            </button>
            <span className="px-2">{page} / {totalPages}</span>
            <button
              onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
              disabled={page === totalPages}
              className="rounded border border-(--color-border) px-2 py-1 disabled:opacity-40"
            >
              ›
            </button>
          </div>
        </div>
      </main>
    </div>
  );
}