import { useState } from "react";
import type { ProjectMember } from "../../types/project";

interface MemberPanelProps {
  members: ProjectMember[];
  loading: boolean;
  onAdd: (email: string) => Promise<void>;
  onRemove: (userId: number) => Promise<void>;
}

export default function MemberPanel({
  members,
  loading,
  onAdd,
  onRemove,
}: MemberPanelProps) {
  const [email, setEmail] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const handleAdd = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email.trim()) return;
    setSubmitting(true);
    try {
      await onAdd(email.trim());
      setEmail("");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="app-card p-5">
      <form onSubmit={handleAdd} className="mb-5 flex gap-2">
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="Enter member's email"
          className="flex-1 rounded-(--radius-button) border border-(--color-border) bg-(--color-bg) px-3 py-2 text-sm text-(--color-text) outline-none focus:border-(--color-primary)"
        />
        <button
          type="submit"
          disabled={submitting}
          className="app-btn-primary flex items-center gap-1 px-4 py-2 text-sm font-medium disabled:opacity-50"
        >
          <i className="bx bx-user-plus"></i>
          Add
        </button>
      </form>

      {loading ? (
        <p className="app-text-muted py-6 text-center text-sm">Loading members...</p>
      ) : members.length === 0 ? (
        <p className="app-text-muted py-6 text-center text-sm">No members yet.</p>
      ) : (
        <ul className="divide-y divide-(--color-border)">
          {members.map((m) => (
            <li
              key={m.projectMemberId}
              className="flex items-center justify-between py-3"
            >
              <div>
                <p className="text-sm font-medium text-(--color-text)">
                  {m.user?.userName || `User #${m.userId}`}
                </p>
                <p className="app-text-muted text-xs">{m.user?.email}</p>
              </div>
              <div className="flex items-center gap-3">
                <span className="rounded-full bg-(--color-bg) px-2 py-0.5 text-xs capitalize text-(--color-muted)">
                  {m.role}
                </span>
                {m.role !== "owner" && (
                  <button
                    onClick={() => onRemove(m.userId)}
                    className="text-red-500 hover:text-red-700"
                    title="Remove member"
                  >
                    <i className="bx bx-user-x text-lg"></i>
                  </button>
                )}
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}