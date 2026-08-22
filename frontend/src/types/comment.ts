export interface Comment {
  commentId: number;
  taskId: number;
  userId: number;
  content: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateCommentPayload {
  content: string;
}