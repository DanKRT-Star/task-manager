import api from "./axios";
import type { Comment, CreateCommentPayload } from "../types/comment";

export const commentService = {
  getTaskComments: async (taskId: number): Promise<Comment[]> => {
    const res = await api.get<{ data: Comment[] }>(`/tasks/${taskId}/comments`);
    return res.data.data;
  },

  createComment: async (
    taskId: number,
    payload: CreateCommentPayload
  ): Promise<Comment> => {
    const res = await api.post<Comment>(`/tasks/${taskId}/comments`, payload);
    return res.data;
  },

  // Route xoá lồng trong task, giống Label chứ không độc lập như Epic/Milestone/Sprint
  deleteComment: async (taskId: number, commentId: number): Promise<void> => {
    await api.delete(`/tasks/${taskId}/comments/${commentId}`);
  },
};