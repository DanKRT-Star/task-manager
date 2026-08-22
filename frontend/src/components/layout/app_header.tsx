import { useState } from "react";
import { toast } from "sonner";
import heroImage from "../../assets/hero.png";
import type { User } from "../../types/user";

interface AppHeaderProps {
  user: User | null;
  onLogout: () => Promise<void>;
}

export default function AppHeader({ user, onLogout }: AppHeaderProps) {
  const [open, setOpen] = useState(false);

  const handleLogout = async () => {
    await onLogout();
    setOpen(false);
    toast.success("Đã đăng xuất");
  };

  return (
    <header className="border-b border-(--color-border) bg-(--color-surface) px-4 py-4 sm:px-6">
      <div className="mx-auto flex max-w-6xl items-center justify-between gap-3">
        <div className="flex items-center gap-2">
          <div className="flex h-8 w-8 items-center justify-center">
            <img src={heroImage} alt="Hero" className="h-8 w-8" />
          </div>
          <span className="text-lg font-semibold text-(--color-text)">
            Task-Manager
          </span>
        </div>

        <div className="relative">
          <button
            onClick={() => setOpen((prev) => !prev)}
            className="flex h-10 w-10 items-center justify-center rounded-full bg-gray-200 text-sm font-semibold text-(--color-text) transition hover:bg-gray-300"
            aria-label="Open user menu"
          >
            {user?.userName?.charAt(0)?.toUpperCase() || "U"}
          </button>

          {open && (
            <div className="absolute right-0 z-10 mt-2 w-64 rounded-(--radius-card) border border-(--color-border) bg-(--color-surface) p-4 shadow-(--shadow-card)">
              <div className="mb-3 border-b border-gray-100 pb-3">
                <p className="text-sm font-semibold text-(--color-text)">
                  {user?.userName || "Unknown user"}
                </p>
                <p className="text-sm text-(--color-muted)">
                  {user?.email || "No email"}
                </p>
              </div>

              <button 
                onClick={handleLogout}
                className="w-full rounded-(--radius-button) bg-(--color-danger) px-3 py-2 text-sm font-medium text-white transition hover:bg-(--color-danger-hover)"
              >
                Log out
              </button>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}