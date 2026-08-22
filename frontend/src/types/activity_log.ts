export type ActivityAction =
  | "created"
  | "status_changed"
  | "assignee_changed"
  | "deadline_changed"
  | "description_updated"
  | "updated"
  | "label_attached"
  | "label_detached";

export interface ActivityLog {
  activityLogId: number;
  taskId: number;
  userId: number;
  action: ActivityAction;
  detail: string;
  createdAt: string;
}