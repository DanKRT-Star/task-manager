import { DragDropContext, Droppable, Draggable, type DropResult } from "@hello-pangea/dnd";
import type { Task, TaskStatus } from "../../types/task";
import TaskCard from "./task_card";

interface TaskKanbanBoardProps {
  tasks: Task[];
  epicTitleById: Map<number, string>;
  milestoneTitleById: Map<number, string>;
  sprintNameById: Map<number, string>;
  onStatusChange: (taskId: number, status: TaskStatus) => void;
  onEdit: (task: Task) => void;
  onDelete: (taskId: number) => void;
  onViewDetail: (task: Task) => void;
  onDropLabel: (taskId: number, labelId: number) => void;
}

const columns: { status: TaskStatus; label: string }[] = [
  { status: "pending", label: "Pending" },
  { status: "in_progress", label: "In Progress" },
  { status: "done", label: "Done" },
];

export default function TaskKanbanBoard({
  tasks,
  epicTitleById,
  milestoneTitleById,
  sprintNameById,
  onStatusChange,
  onEdit,
  onDelete,
  onViewDetail,
  onDropLabel,
}: TaskKanbanBoardProps) {
  const handleDragEnd = (result: DropResult) => {
    const { source, destination, draggableId } = result;
    if (!destination) return; // thả ra ngoài mọi Droppable -> huỷ
    if (destination.droppableId === source.droppableId) return; // thả lại đúng cột cũ

    onStatusChange(Number(draggableId), destination.droppableId as TaskStatus);
  };

  return (
    <DragDropContext onDragEnd={handleDragEnd}>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
        {columns.map((col) => {
          const colTasks = tasks.filter((t) => t.status === col.status);
          return (
            <Droppable droppableId={col.status} key={col.status}>
              {(provided, snapshot) => (
                <div
                  ref={provided.innerRef}
                  {...provided.droppableProps}
                  className={`flex min-h-24 flex-col gap-3 rounded-(--radius-button) p-1 transition-colors ${
                    snapshot.isDraggingOver ? "bg-(--color-bg)" : ""
                  }`}
                >
                  <div className="flex items-center justify-between px-1">
                    <h3 className="text-sm font-semibold text-(--color-text)">
                      {col.label}
                    </h3>
                    <span className="app-text-muted text-xs">{colTasks.length}</span>
                  </div>

                  {colTasks.length === 0 && !snapshot.isDraggingOver && (
                    <div className="app-card py-8 text-center">
                      <p className="app-text-muted text-xs">No tasks</p>
                    </div>
                  )}

                  {colTasks.map((task, index) => (
                    <Draggable
                      key={task.taskId}
                      draggableId={String(task.taskId)}
                      index={index}
                    >
                      {(dragProvided, dragSnapshot) => (
                        <TaskCard
                          task={task}
                          epicTitle={task.epicId ? epicTitleById.get(task.epicId) : undefined}
                          milestoneTitle={
                            task.milestoneId
                              ? milestoneTitleById.get(task.milestoneId)
                              : undefined
                          }
                          sprintName={
                            task.sprintId ? sprintNameById.get(task.sprintId) : undefined
                          }
                          onStatusChange={onStatusChange}
                          onEdit={onEdit}
                          onDelete={onDelete}
                          onViewDetail={onViewDetail}
                          onDropLabel={onDropLabel}
                          dragProvided={dragProvided}
                          isDragging={dragSnapshot.isDragging}
                        />
                      )}
                    </Draggable>
                  ))}

                  {provided.placeholder}
                </div>
              )}
            </Droppable>
          );
        })}
      </div>
    </DragDropContext>
  );
}