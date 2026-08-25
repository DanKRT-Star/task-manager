import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  Calendar,
  dateFnsLocalizer,
  Views,
  type View,
  type Event as RBCEvent,
  type EventProps,
} from "react-big-calendar";
import "react-big-calendar/lib/css/react-big-calendar.css";
import { format } from "date-fns/format";
import { parse } from "date-fns/parse";
import { startOfWeek } from "date-fns/startOfWeek";
import { getDay } from "date-fns/getDay";
import { enUS } from "date-fns/locale/en-US";
import { useTask } from "../hooks/use_task";
import TaskFormPanel from "../components/task/task_form_panel";
import type { Task } from "../types/task";

const locales = { "en-US": enUS };

const localizer = dateFnsLocalizer({
  format,
  parse,
  startOfWeek: () => startOfWeek(new Date(), { locale: enUS }),
  getDay,
  locales,
});

interface TaskEvent extends RBCEvent {
  task: Task;
}

export default function CalendarPage() {
  const { tasks, loading, fetchTasks, updateTask, createTask } = useTask();
  const navigate = useNavigate();

  const [panelOpen, setPanelOpen] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);

  const [currentDate, setCurrentDate] = useState(new Date());
  const [currentView, setCurrentView] = useState<View>(Views.MONTH);

  useEffect(() => {
    fetchTasks({ limit: 100, sort: "deadline_asc" });
  }, [fetchTasks]);

  const events: TaskEvent[] = useMemo(
    () =>
      tasks
        .filter((t) => !!t.deadline)
        .map((t) => {
          const date = new Date(t.deadline);
          return {
            title: t.title,
            start: date,
            end: date,
            allDay: true,
            task: t,
          };
        }),
    [tasks]
  );

  const handleSelectEvent = (event: TaskEvent) => {
    const { task } = event;
    if (task.project) {
      navigate(`/projects/${task.project.projectId}`);
      return;
    }
    setEditingTask(task);
    setPanelOpen(true);
  };

  const handleSelectSlot = ({ start }: { start: Date }) => {
    setEditingTask({
      taskId: 0,
      title: "",
      description: "",
      status: "pending",
      userId: 0,
      deadline: start.toISOString(),
      createdAt: "",
      updatedAt: "",
    } as Task);
    setPanelOpen(true);
  };

  const handleSubmit = async (data: {
    title: string;
    description?: string;
    status: Task["status"];
    deadline: string;
  }) => {
    try {
      const payload = { ...data, deadline: new Date(data.deadline).toISOString() };
      if (editingTask && editingTask.taskId !== 0) {
        await updateTask(editingTask.taskId, payload);
      } else {
        await createTask(payload);
      }
      fetchTasks({ limit: 100, sort: "deadline_asc" });
    } catch {
      // toast đã xử lý trong hook
    }
  };

  const eventPropGetter = (event: TaskEvent) => {
    const { task } = event;
    let backgroundColor = "var(--color-primary)"; // task cá nhân, pending
    if (task.status === "done") backgroundColor = "#22c55e";
    else if (task.status === "in_progress") backgroundColor = "#3b82f6";
    else if (task.project) backgroundColor = "#a855f7"; // task project (mình là assignee), pending

    return {
      style: {
        backgroundColor,
        borderRadius: "6px",
        border: "none",
        fontSize: "0.75rem",
      },
    };
  };

  const EventContent = ({ event }: EventProps<TaskEvent>) => (
    <span title={event.task.project ? `Project: ${event.task.project.name}` : undefined}>
      {event.task.project && <i className="bx bx-briefcase mr-1"></i>}
      {event.title}
    </span>
  );

  return (
    <div className="mx-auto max-w-6xl p-4 sm:p-6">
      <div className="mb-4 flex items-center justify-between">
        <h1 className="app-panel-title">Calendar</h1>
        <div className="flex items-center gap-3 text-xs">
          <span className="flex items-center gap-1">
            <span className="h-2.5 w-2.5 rounded-full" style={{ backgroundColor: "var(--color-primary)" }}></span>
            Personal
          </span>
          <span className="flex items-center gap-1">
            <span className="h-2.5 w-2.5 rounded-full bg-purple-500"></span>
            From project
          </span>
          <span className="flex items-center gap-1">
            <span className="h-2.5 w-2.5 rounded-full bg-green-500"></span>
            Done
          </span>
        </div>
      </div>

      <div className="app-card p-4">
        {loading && tasks.length === 0 ? (
          <div className="py-20 text-center text-(--color-muted)">Loading...</div>
        ) : (
          <Calendar
            localizer={localizer}
            events={events}
            startAccessor="start"
            endAccessor="end"
            style={{ height: 700 }}
            selectable
            date={currentDate}
            onNavigate={setCurrentDate}
            view={currentView}
            onView={setCurrentView}
            onSelectEvent={handleSelectEvent}
            onSelectSlot={handleSelectSlot}
            eventPropGetter={eventPropGetter}
            components={{ event: EventContent }}
            popup
          />
        )}
      </div>

      <TaskFormPanel
        open={panelOpen}
        onClose={() => setPanelOpen(false)}
        onSubmit={handleSubmit}
        editingTask={editingTask && editingTask.taskId !== 0 ? editingTask : null}
      />
    </div>
  );
}