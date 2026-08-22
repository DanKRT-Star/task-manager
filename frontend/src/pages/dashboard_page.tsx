import { useAuth } from "../hooks/use_auth";

export default function DashboardPage() {
  const { user } = useAuth();

  return (
    <div className="mx-auto max-w-6xl p-4 sm:p-6">
      <h1 className="app-panel-title mb-6">
        Welcome, {user?.userName || "User"}!
      </h1>
      
      <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
        <div className="app-card col-span-1 p-6 md:col-span-2">
          <h2 className="mb-4 text-lg font-bold text-(--color-text)">
            My Tasks
          </h2>
          <p className="app-text-muted text-sm">
            Tasks to complete today will appear here...
          </p>
        </div>
        
        <div className="flex flex-col gap-6">
          <div className="app-card p-6">
            <h2 className="mb-4 text-lg font-bold text-(--color-text)">
              Recent Projects
            </h2>
            <p className="app-text-muted text-sm">
              New Website, Marketing Campaign...
            </p>
          </div>
          
          <div className="app-card flex-1 p-6">
            <h2 className="mb-4 text-lg font-bold text-(--color-text)">
              Recent Activity
            </h2>
            <p className="app-text-muted text-sm">
              Team activity stream...
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}