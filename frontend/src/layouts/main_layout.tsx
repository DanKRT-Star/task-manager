import { Outlet } from "react-router-dom";
import { useAuth } from "../hooks/use_auth";
import AppHeader from "../components/layout/app_header";
import Sidebar from "../components/layout/sidebar";

export default function MainLayout() {
  const { user, logout } = useAuth();

  return (
    <div className="app-shell flex h-screen w-full overflow-hidden">
      <Sidebar />

      <div className="flex flex-1 flex-col overflow-hidden">
        <AppHeader user={user} onLogout={logout} />
        
        <main className="flex-1 overflow-y-auto">
          <Outlet />
        </main>
      </div>
    </div>
  );
}