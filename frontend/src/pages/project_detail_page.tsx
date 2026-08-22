import { useEffect, useMemo, useState } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import { useProject } from "../hooks/use_project";
import { useEpic } from "../hooks/use_epic";
import { useMilestone } from "../hooks/use_milestone";
import { useSprint } from "../hooks/use_sprint";
import { useLabel } from "../hooks/use_label";
import { useComment } from "../hooks/use_comment";
import { useActivityLog } from "../hooks/use_activity_log";
import { useAuth } from "../hooks/use_auth";
import { useTask } from "../hooks/use_task";
import ProjectFormPanel from "../components/project/project_form_panel";
import EpicFormPanel from "../components/epic/epic_form_panel";
import EpicsTab from "../components/epic/epic_tab";
import MilestoneFormPanel from "../components/milestone/milestone_form_panel";
import MilestonesTab from "../components/milestone/milestone_tab";
import SprintFormPanel from "../components/sprint/sprint_form_panel";
import SprintsTab from "../components/sprint/sprint_tab";
import LabelManagerPanel from "../components/label/label_manager_panel";
import MemberPanel from "../components/project/member_panel";
import TaskFilterBar, { type TaskView } from "../components/task/task_filter_bar";
import TaskKanbanBoard from "../components/task/task_karban_board";
import TaskListView from "../components/task/task_list_view";
import ProjectTaskFormPanel from "../components/task/project_task_form_panel";
import TaskDetailPanel from "../components/task/task_detail_panel";
import type { Epic } from "../types/epic";
import type { Milestone } from "../types/milestone";
import type { Sprint, SprintStatus, CreateSprintPayload } from "../types/sprint";
import type { Task, TaskStatus, CreateTaskPayload } from "../types/task";
import type { ProjectFormData } from "../lib/validation";

type Tab = "tasks" | "epics" | "milestones" | "sprints" | "labels" | "members";

export default function ProjectDetailPage() {
  const { id } = useParams<{ id: string }>();
  const projectId = Number(id);
  const navigate = useNavigate();

  const {
    currentProject,
    members,
    loading: projectLoading,
    fetchProject,
    updateProject,
    deleteProject,
    fetchMembers,
    addMember,
    removeMember,
  } = useProject();

  const {
    epics,
    loading: epicsLoading,
    fetchEpics,
    createEpic,
    updateEpic,
    deleteEpic,
  } = useEpic();

  const {
    milestones,
    loading: milestonesLoading,
    fetchMilestones,
    createMilestone,
    updateMilestone,
    deleteMilestone,
  } = useMilestone();

  const {
    sprints,
    loading: sprintsLoading,
    fetchSprints,
    createSprint,
    updateSprint,
    deleteSprint,
  } = useSprint();

  const {
    labels,
    loading: labelsLoading,
    fetchLabels,
    createLabel,
    deleteLabel,
    attachLabel,
    detachLabel,
  } = useLabel();

  const {
    tasks,
    loading: tasksLoading,
    fetchProjectTasks,
    createTask,
    updateTask,
    deleteTask,
  } = useTask();

  const {
    comments,
    loading: commentsLoading,
    fetchComments,
    createComment,
    deleteComment,
  } = useComment();

  const {
    logs: activityLogs,
    loading: activityLoading,
    fetchTaskActivity,
  } = useActivityLog();

  const { user } = useAuth();

  const [tab, setTab] = useState<Tab>("tasks");

  const [projectPanelOpen, setProjectPanelOpen] = useState(false);

  const [epicPanelOpen, setEpicPanelOpen] = useState(false);
  const [editingEpic, setEditingEpic] = useState<Epic | null>(null);

  const [milestonePanelOpen, setMilestonePanelOpen] = useState(false);
  const [editingMilestone, setEditingMilestone] = useState<Milestone | null>(null);

  const [sprintPanelOpen, setSprintPanelOpen] = useState(false);
  const [editingSprint, setEditingSprint] = useState<Sprint | null>(null);

  const [taskPanelOpen, setTaskPanelOpen] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);

  const [detailPanelOpen, setDetailPanelOpen] = useState(false);
  const [viewingTask, setViewingTask] = useState<Task | null>(null);

  // Bộ lọc + view cho danh sách Task — Epic/Milestone/Sprint/Label đều là filter chồng lên Task
  const [taskView, setTaskView] = useState<TaskView>("kanban");
  const [epicFilter, setEpicFilter] = useState<number | null>(null);
  const [milestoneFilter, setMilestoneFilter] = useState<number | null>(null);
  const [sprintFilter, setSprintFilter] = useState<number | null>(null);
  const [labelFilter, setLabelFilter] = useState<number[]>([]);
  const [statusFilter, setStatusFilter] = useState<TaskStatus | null>(null);

  useEffect(() => {
    if (!projectId) return;
    fetchProject(projectId).catch(() => {});
    fetchEpics(projectId);
    fetchMilestones(projectId);
    fetchSprints(projectId);
    fetchLabels(projectId);
    fetchMembers(projectId);
    fetchProjectTasks(projectId);
  }, [
    projectId,
    fetchProject,
    fetchEpics,
    fetchMilestones,
    fetchSprints,
    fetchLabels,
    fetchMembers,
    fetchProjectTasks,
  ]);

  const epicTitleById = useMemo(
    () => new Map(epics.map((e) => [e.epicId, e.title])),
    [epics]
  );
  const milestoneTitleById = useMemo(
    () => new Map(milestones.map((m) => [m.milestoneId, m.title])),
    [milestones]
  );
  const sprintNameById = useMemo(
    () => new Map(sprints.map((s) => [s.sprintId, s.name])),
    [sprints]
  );
  const memberNameById = useMemo(
    () =>
      new Map(
        members.map((m) => [m.userId, m.user?.userName || `User #${m.userId}`])
      ),
    [members]
  );

  const filteredTasks = useMemo(() => {
    return tasks.filter((t) => {
      if (epicFilter !== null && t.epicId !== epicFilter) return false;
      if (milestoneFilter !== null && t.milestoneId !== milestoneFilter) return false;
      if (sprintFilter !== null && t.sprintId !== sprintFilter) return false;
      if (statusFilter !== null && t.status !== statusFilter) return false;
      if (labelFilter.length > 0) {
        const taskLabelIds = t.labels?.map((l) => l.labelId) || [];
        if (!labelFilter.some((id) => taskLabelIds.includes(id))) return false;
      }
      return true;
    });
  }, [tasks, epicFilter, milestoneFilter, sprintFilter, labelFilter, statusFilter]);

  // ---- Project actions ----
  const handleUpdateProject = async (data: ProjectFormData) => {
    try {
      const payload = {
        ...data,
        deadline: data.deadline ? new Date(data.deadline).toISOString() : undefined,
      };
      await updateProject(projectId, payload);
    } catch {
      // Toast is already handled inside the hook
    }
  };

  const handleDeleteProject = async () => {
    if (!confirm("Are you sure you want to delete this project?")) return;
    try {
      await deleteProject(projectId);
      navigate("/projects");
    } catch {
      // Toast is already handled inside the hook
    }
  };

  // ---- Epic actions ----
  const handleCreateEpic = () => {
    setEditingEpic(null);
    setEpicPanelOpen(true);
  };

  const handleEditEpic = (epic: Epic) => {
    setEditingEpic(epic);
    setEpicPanelOpen(true);
  };

  const handleEpicSubmit = async (data: { title: string; description?: string }) => {
    try {
      if (editingEpic) {
        await updateEpic(editingEpic.epicId, data);
      } else {
        await createEpic(projectId, data);
      }
    } catch {
      // Toast is already handled inside the hook
    }
  };

  const handleDeleteEpic = async (epicId: number) => {
    if (!confirm("Delete this epic? Tasks in it will remain but lose the epic link.")) return;
    try {
      await deleteEpic(epicId);
      if (epicFilter === epicId) setEpicFilter(null);
    } catch {
      // Toast is already handled inside the hook
    }
  };

  const handleViewEpicTasks = (epicId: number) => {
    setEpicFilter(epicId);
    setTab("tasks");
  };

  // ---- Milestone actions ----
  const handleCreateMilestone = () => {
    setEditingMilestone(null);
    setMilestonePanelOpen(true);
  };

  const handleEditMilestone = (milestone: Milestone) => {
    setEditingMilestone(milestone);
    setMilestonePanelOpen(true);
  };

  const handleMilestoneSubmit = async (data: {
    title: string;
    description?: string;
    dueDate?: string;
  }) => {
    try {
      const payload = {
        ...data,
        dueDate: data.dueDate ? new Date(data.dueDate).toISOString() : undefined,
      };
      if (editingMilestone) {
        await updateMilestone(editingMilestone.milestoneId, payload);
      } else {
        await createMilestone(projectId, payload);
      }
    } catch {
      // Toast is already handled inside the hook
    }
  };

  const handleDeleteMilestone = async (milestoneId: number) => {
    if (!confirm("Delete this milestone? Tasks in it will remain but lose the milestone link.")) return;
    try {
      await deleteMilestone(milestoneId);
      if (milestoneFilter === milestoneId) setMilestoneFilter(null);
    } catch {
      // Toast is already handled inside the hook
    }
  };

  const handleViewMilestoneTasks = (milestoneId: number) => {
    setMilestoneFilter(milestoneId);
    setTab("tasks");
  };

  // ---- Sprint actions ----
  const handleCreateSprint = () => {
    setEditingSprint(null);
    setSprintPanelOpen(true);
  };

  const handleEditSprint = (sprint: Sprint) => {
    setEditingSprint(sprint);
    setSprintPanelOpen(true);
  };

  const handleSprintSubmit = async (
    data: CreateSprintPayload & { status?: SprintStatus }
  ) => {
    try {
      if (editingSprint) {
        await updateSprint(editingSprint.sprintId, data);
      } else {
        await createSprint(projectId, data);
      }
    } catch {
      // Toast is already handled inside the hook
    }
  };

  const handleDeleteSprint = async (sprintId: number) => {
    if (!confirm("Delete this sprint? Tasks in it will remain but lose the sprint link.")) return;
    try {
      await deleteSprint(sprintId);
      if (sprintFilter === sprintId) setSprintFilter(null);
    } catch {
      // Toast is already handled inside the hook
    }
  };

  const handleViewSprintTasks = (sprintId: number) => {
    setSprintFilter(sprintId);
    setTab("tasks");
  };

  // ---- Label actions ----
  const handleCreateLabel = async (payload: { name: string; color?: string }) => {
    await createLabel(projectId, payload);
  };

  const handleDeleteLabel = async (labelId: number) => {
    await deleteLabel(projectId, labelId);
    setLabelFilter((prev) => prev.filter((id) => id !== labelId));
  };

  // ---- Member actions ----
  const handleAddMember = async (email: string) => {
    await addMember(projectId, email);
  };

  const handleRemoveMember = async (userId: number) => {
    if (!confirm("Remove this member from the project?")) return;
    await removeMember(projectId, userId);
  };

  // ---- Task actions ----
  const handleCreateTask = () => {
    setEditingTask(null);
    setTaskPanelOpen(true);
  };

  const handleEditTask = (task: Task) => {
    setEditingTask(task);
    setTaskPanelOpen(true);
  };

  const handleTaskSubmit = async (data: CreateTaskPayload, selectedLabelIds: number[]) => {
    try {
      let savedTask: Task;
      if (editingTask) {
        savedTask = await updateTask(editingTask.taskId, data);
      } else {
        savedTask = await createTask({ ...data, projectId });
      }

      const currentLabelIds = editingTask?.labels?.map((l) => l.labelId) || [];
      const toAttach = selectedLabelIds.filter((id) => !currentLabelIds.includes(id));
      const toDetach = currentLabelIds.filter((id) => !selectedLabelIds.includes(id));

      await Promise.all([
        ...toAttach.map((labelId) => attachLabel(savedTask.taskId, labelId)),
        ...toDetach.map((labelId) => detachLabel(savedTask.taskId, labelId)),
      ]);

      if (toAttach.length > 0 || toDetach.length > 0) {
        await fetchProjectTasks(projectId);
      }
    } catch {
      // Toast is already handled inside the hooks
    }
  };

  const handleTaskStatusChange = async (taskId: number, status: TaskStatus) => {
    try {
      await updateTask(taskId, { status });
    } catch {
      // Toast is already handled inside the hook
    }
  };

  const handleDeleteTask = async (taskId: number) => {
    if (!confirm("Delete this task?")) return;
    try {
      await deleteTask(taskId);
    } catch {
      // Toast is already handled inside the hook
    }
  };

  const toggleLabelFilter = (labelId: number) => {
    setLabelFilter((prev) =>
      prev.includes(labelId) ? prev.filter((id) => id !== labelId) : [...prev, labelId]
    );
  };

  const handleDropLabelOnTask = async (taskId: number, labelId: number) => {
    const task = tasks.find((t) => t.taskId === taskId);
    if (task?.labels?.some((l) => l.labelId === labelId)) return;
    try {
      await attachLabel(taskId, labelId);
      await fetchProjectTasks(projectId);
    } catch {
      // Toast is already handled inside the hook
    }
  };

  // ---- Task detail (comments + activity) ----
  const handleViewTaskDetail = (task: Task) => {
    setViewingTask(task);
    setDetailPanelOpen(true);
  };

  const handleCreateComment = async (taskId: number, content: string) => {
    try {
      await createComment(taskId, { content });
    } catch {
      // Toast is already handled inside the hook
    }
  };

  const handleDeleteComment = async (taskId: number, commentId: number) => {
    if (!confirm("Delete this comment?")) return;
    try {
      await deleteComment(taskId, commentId);
    } catch {
      // Toast is already handled inside the hook
    }
  };

  if (!projectId) {
    return (
      <div className="mx-auto max-w-6xl p-6">
        <p className="app-text-muted">Invalid project id.</p>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-6xl p-4 sm:p-6">
      <Link to="/projects" className="app-text-muted mb-4 inline-flex items-center gap-1 text-sm hover:text-(--color-text)">
        <i className="bx bx-arrow-back"></i>
        Back to Projects
      </Link>

      {projectLoading && !currentProject ? (
        <div className="app-card py-12 text-center">
          <p className="app-text-muted">Loading project...</p>
        </div>
      ) : !currentProject ? (
        <div className="app-card py-12 text-center">
          <p className="app-text-muted">Project not found.</p>
        </div>
      ) : (
        <>
          <div className="app-card mb-6 flex flex-col gap-4 p-5 sm:flex-row sm:items-start sm:justify-between">
            <div>
              <h1 className="app-panel-title mb-1 text-2xl">{currentProject.name}</h1>
              <p className="app-text-muted max-w-2xl text-sm">
                {currentProject.description || "No description provided."}
              </p>
              <div className="app-text-muted mt-3 flex items-center gap-1 text-xs">
                <i className="bx bx-time-five text-base"></i>
                {currentProject.deadline
                  ? currentProject.deadline.slice(0, 10)
                  : "No deadline"}
              </div>
            </div>

            <div className="flex shrink-0 gap-2">
              <button
                onClick={() => setProjectPanelOpen(true)}
                className="app-btn-secondary flex items-center gap-1 px-3 py-2 text-sm font-medium"
              >
                <i className="bx bx-edit"></i>
                Edit
              </button>
              <button
                onClick={handleDeleteProject}
                className="app-btn-secondary flex items-center gap-1 px-3 py-2 text-sm font-medium text-red-500 hover:text-red-700"
              >
                <i className="bx bx-trash"></i>
                Delete
              </button>
            </div>
          </div>

          <div className="mb-5 flex gap-2 overflow-x-auto border-b border-(--color-border)">
            {(
              [
                { key: "tasks", label: "Tasks", count: tasks.length },
                { key: "epics", label: "Epics", count: epics.length },
                { key: "milestones", label: "Milestones", count: milestones.length },
                { key: "sprints", label: "Sprints", count: sprints.length },
                { key: "labels", label: "Labels", count: labels.length },
                { key: "members", label: "Members", count: members.length },
              ] as { key: Tab; label: string; count: number }[]
            ).map((t) => (
              <button
                key={t.key}
                onClick={() => setTab(t.key)}
                className={`-mb-px flex shrink-0 items-center gap-1.5 border-b-2 px-3 py-2 text-sm font-medium ${
                  tab === t.key
                    ? "border-(--color-primary) text-(--color-primary)"
                    : "border-transparent app-text-muted hover:text-(--color-text)"
                }`}
              >
                {t.label}
                <span className="rounded-full bg-(--color-bg) px-1.5 text-xs">
                  {t.count}
                </span>
              </button>
            ))}
          </div>

          {tab === "tasks" && (
            <div>
              <TaskFilterBar
                epics={epics}
                milestones={milestones}
                sprints={sprints}
                labels={labels}
                epicFilter={epicFilter}
                milestoneFilter={milestoneFilter}
                sprintFilter={sprintFilter}
                labelFilter={labelFilter}
                statusFilter={statusFilter}
                view={taskView}
                onEpicFilterChange={setEpicFilter}
                onMilestoneFilterChange={setMilestoneFilter}
                onSprintFilterChange={setSprintFilter}
                onLabelFilterToggle={toggleLabelFilter}
                onStatusFilterChange={setStatusFilter}
                onViewChange={setTaskView}
                onNewTask={handleCreateTask}
              />

              {tasksLoading ? (
                <div className="app-card py-12 text-center">
                  <p className="app-text-muted">Loading tasks...</p>
                </div>
              ) : taskView === "kanban" ? (
                <TaskKanbanBoard
                  tasks={filteredTasks}
                  epicTitleById={epicTitleById}
                  milestoneTitleById={milestoneTitleById}
                  sprintNameById={sprintNameById}
                  onStatusChange={handleTaskStatusChange}
                  onEdit={handleEditTask}
                  onDelete={handleDeleteTask}
                  onViewDetail={handleViewTaskDetail}
                  onDropLabel={handleDropLabelOnTask}
                />
              ) : (
                <TaskListView
                  tasks={filteredTasks}
                  epicTitleById={epicTitleById}
                  milestoneTitleById={milestoneTitleById}
                  sprintNameById={sprintNameById}
                  onStatusChange={handleTaskStatusChange}
                  onEdit={handleEditTask}
                  onDelete={handleDeleteTask}
                  onViewDetail={handleViewTaskDetail}
                />
              )}
            </div>
          )}

          {tab === "epics" && (
            <EpicsTab
              epics={epics}
              loading={epicsLoading}
              onCreate={handleCreateEpic}
              onEdit={handleEditEpic}
              onDelete={handleDeleteEpic}
              onViewTasks={handleViewEpicTasks}
            />
          )}

          {tab === "milestones" && (
            <MilestonesTab
              milestones={milestones}
              loading={milestonesLoading}
              onCreate={handleCreateMilestone}
              onEdit={handleEditMilestone}
              onDelete={handleDeleteMilestone}
              onViewTasks={handleViewMilestoneTasks}
            />
          )}

          {tab === "sprints" && (
            <SprintsTab
              sprints={sprints}
              loading={sprintsLoading}
              onCreate={handleCreateSprint}
              onEdit={handleEditSprint}
              onDelete={handleDeleteSprint}
              onViewTasks={handleViewSprintTasks}
            />
          )}

          {tab === "labels" && (
            <LabelManagerPanel
              labels={labels}
              loading={labelsLoading}
              onCreate={handleCreateLabel}
              onDelete={handleDeleteLabel}
            />
          )}

          {tab === "members" && (
            <MemberPanel
              members={members}
              loading={projectLoading}
              onAdd={handleAddMember}
              onRemove={handleRemoveMember}
            />
          )}

          <ProjectFormPanel
            open={projectPanelOpen}
            onClose={() => setProjectPanelOpen(false)}
            onSubmit={handleUpdateProject}
            editingProject={currentProject}
          />

          <EpicFormPanel
            open={epicPanelOpen}
            onClose={() => setEpicPanelOpen(false)}
            onSubmit={handleEpicSubmit}
            editingEpic={editingEpic}
          />

          <MilestoneFormPanel
            open={milestonePanelOpen}
            onClose={() => setMilestonePanelOpen(false)}
            onSubmit={handleMilestoneSubmit}
            editingMilestone={editingMilestone}
          />

          <SprintFormPanel
            open={sprintPanelOpen}
            onClose={() => setSprintPanelOpen(false)}
            onSubmit={handleSprintSubmit}
            editingSprint={editingSprint}
          />

          <ProjectTaskFormPanel
            open={taskPanelOpen}
            onClose={() => setTaskPanelOpen(false)}
            onSubmit={handleTaskSubmit}
            editingTask={editingTask}
            epics={epics}
            milestones={milestones}
            sprints={sprints}
            labels={labels}
            members={members}
            defaultEpicId={epicFilter}
            defaultMilestoneId={milestoneFilter}
            defaultSprintId={sprintFilter}
          />

          <TaskDetailPanel
            key={detailPanelOpen && viewingTask ? viewingTask.taskId : "closed"}
            open={detailPanelOpen}
            onClose={() => setDetailPanelOpen(false)}
            task={viewingTask}
            currentUserId={user?.userId}
            memberNameById={memberNameById}
            comments={comments}
            commentsLoading={commentsLoading}
            activityLogs={activityLogs}
            activityLoading={activityLoading}
            onFetchComments={fetchComments}
            onFetchActivity={fetchTaskActivity}
            onCreateComment={handleCreateComment}
            onDeleteComment={handleDeleteComment}
          />
        </>
      )}
    </div>
  );
}