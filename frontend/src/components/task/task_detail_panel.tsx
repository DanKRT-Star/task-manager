import { useEffect, useState } from "react";
import type { Task } from "../../types/task";
import type { Comment } from "../../types/comment";
import type { ActivityLog, ActivityAction } from "../../types/activity_log";

interface TaskDetailPanelProps {
  open: boolean;
  onClose: () => void;
  task: Task | null;
  currentUserId?: number;
  memberNameById: Map<number, string>;
  comments: Comment[];
  commentsLoading: boolean;
  activityLogs: ActivityLog[];
  activityLoading: boolean;
  onFetchComments: (taskId: number) => void;
  onFetchActivity: (taskId: number) => void;
  onCreateComment: (taskId: number, content: string) => Promise<void>;
  onDeleteComment: (taskId: number, commentId: number) => Promise<void>;
}

type DetailTab = "comments" | "activity";

const actionLabel: Record<ActivityAction, string> = {
  created: "created the task",
  status_changed: "changed status",
  assignee_changed: "changed assignee",
  deadline_changed: "changed the deadline",
  description_updated: "updated the description",
  updated: "updated the task",
  label_attached: "attached a label",
  label_detached: "removed a label",
};

export default function TaskDetailPanel({
  open,
  onClose,
  task,
  currentUserId,
  memberNameById,
  comments,
  commentsLoading,
  activityLogs,
  activityLoading,
  onFetchComments,
  onFetchActivity,
  onCreateComment,
  onDeleteComment,
}: TaskDetailPanelProps) {
  const [tab, setTab] = useState<DetailTab>("comments");
  const [content, setContent] = useState("");
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (open && task) {
      onFetchComments(task.taskId);
      onFetchActivity(task.taskId);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, task?.taskId]);

  if (!open || !task) return null;

  const resolveName = (userId: number) =>
    userId === currentUserId ? "You" : memberNameById.get(userId) || `User #${userId}`;

  const handleSubmitComment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!content.trim()) return;
    setSubmitting(true);
    try {
      await onCreateComment(task.taskId, content.trim());
      setContent("");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <>
      <div
        className="fixed inset-0 z-10 bg-black/20 transition-opacity"
        onClick={onClose}
      />

      <div className="fixed right-0 top-0 z-20 flex h-full w-96 flex-col border-l border-(--color-border) bg-(--color-surface) p-6 shadow-lg transition-transform">
        <div className="mb-4 flex items-start justify-between gap-3">
          <h2 className="app-panel-title line-clamp-2 text-lg">{task.title}</h2>
          <button
            onClick={onClose}
            className="shrink-0 text-(--color-muted) hover:text-(--color-text)"
          >
            <i className="bx bx-x text-2xl"></i>
          </button>
        </div>

        <div className="mb-4 flex gap-2 border-b border-(--color-border)">
          <button
            onClick={() => setTab("comments")}
            className={`-mb-px border-b-2 px-3 py-2 text-sm font-medium ${
              tab === "comments"
                ? "border-(--color-primary) text-(--color-primary)"
                : "border-transparent app-text-muted hover:text-(--color-text)"
            }`}
          >
            Comments
          </button>
          <button
            onClick={() => setTab("activity")}
            className={`-mb-px border-b-2 px-3 py-2 text-sm font-medium ${
              tab === "activity"
                ? "border-(--color-primary) text-(--color-primary)"
                : "border-transparent app-text-muted hover:text-(--color-text)"
            }`}
          >
            Activity
          </button>
        </div>

        <div className="flex-1 overflow-y-auto">
          {tab === "comments" ? (
            commentsLoading ? (
              <p className="app-text-muted py-6 text-center text-sm">Loading comments...</p>
            ) : comments.length === 0 ? (
              <p className="app-text-muted py-6 text-center text-sm">No comments yet.</p>
            ) : (
              <ul className="space-y-3">
                {comments.map((c) => (
                  <li key={c.commentId} className="app-card p-3">
                    <div className="mb-1 flex items-center justify-between">
                      <span className="text-xs font-medium text-(--color-text)">
                        {resolveName(c.userId)}
                        <span className="app-text-muted ml-2 font-normal">
                          {new Date(c.createdAt).toLocaleString()}
                        </span>
                      </span>
                      {currentUserId === c.userId && (
                        <button
                          onClick={() => onDeleteComment(task.taskId, c.commentId)}
                          className="text-red-500 hover:text-red-700"
                        >
                          <i className="bx bx-trash text-sm"></i>
                        </button>
                      )}
                    </div>
                    <p className="text-sm text-(--color-text)">{c.content}</p>
                  </li>
                ))}
              </ul>
            )
          ) : activityLoading ? (
            <p className="app-text-muted py-6 text-center text-sm">Loading activity...</p>
          ) : activityLogs.length === 0 ? (
            <p className="app-text-muted py-6 text-center text-sm">No activity yet.</p>
          ) : (
            <ul className="space-y-3">
              {activityLogs.map((log) => (
                <li key={log.activityLogId} className="flex gap-2 text-sm">
                  <i className="bx bx-history mt-0.5 text-(--color-muted)"></i>
                  <div>
                    <p className="text-(--color-text)">
                      <span className="font-medium">{resolveName(log.userId)}</span>{" "}
                      {actionLabel[log.action] || log.action}
                      {log.detail ? ` — ${log.detail}` : ""}
                    </p>
                    <span className="app-text-muted text-xs">
                      {new Date(log.createdAt).toLocaleString()}
                    </span>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </div>

        {tab === "comments" && (
          <form onSubmit={handleSubmitComment} className="mt-4 flex gap-2 border-t border-(--color-border) pt-4">
            <input
              type="text"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              placeholder="Write a comment..."
              className="flex-1 rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
            />
            <button
              type="submit"
              disabled={submitting}
              className="app-btn-primary px-3 py-2 text-sm font-medium disabled:opacity-50"
            >
              <i className="bx bx-send"></i>
            </button>
          </form>
        )}
      </div>
    </>
  );
}