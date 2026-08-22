import { Link, useLocation } from "react-router-dom";

export default function Sidebar() {
  const location = useLocation();
  
  const navItems = [
    { path: "/", label: "Dashboard", icon: "bx-home" },
    { path: "/tasks", label: "My Tasks", icon: "bx-check-square" },
    { path: "/projects", label: "Projects", icon: "bx-folder" },
  ];

  return (
    <aside className="hidden w-64 flex-col border-r border-(--color-border) bg-(--color-surface) md:flex">
      <div className="flex items-center gap-2 border-b border-(--color-border) p-4 text-xl font-bold text-(--color-text)">
        <i className="bx bx-check-shield text-(--color-primary)"></i>
        ProTask
      </div>
      
      <nav className="flex-1 space-y-2 p-4">
        {navItems.map((item) => {
          const isActive =
            location.pathname === item.path ||
            (item.path !== "/" && location.pathname.startsWith(item.path));
          
          return (
            <Link
              key={item.path}
              to={item.path}
              className={`flex items-center gap-3 rounded-(--radius-button) px-4 py-2.5 transition-colors ${
                isActive
                  ? "bg-(--color-primary) font-medium text-white"
                  : "text-(--color-muted) hover:bg-(--color-bg) hover:text-(--color-text)"
              }`}
            >
              <i className={`bx ${item.icon} text-xl`}></i>
              {item.label}
            </Link>
          );
        })}
      </nav>
    </aside>
  );
}