import { useState, useCallback } from "react";
import { toast } from "sonner";
import { commentService } from "../services/comment_service";
import type { Comment, CreateCommentPayload } from "../types/comment";
import { getErrorMessage } from "../lib/error";

export function useComment() {
  const [comments, setComments] = useState<Comment[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchComments = useCallback(async (taskId: number) => {
    setLoading(true);
    try {
      const data = await commentService.getTaskComments(taskId);
      setComments(data);
    } catch (err) {
      toast.error(getErrorMessage(err, "Không tải được bình luận"));
    } finally {
      setLoading(false);
    }
  }, []);

  const createComment = async (taskId: number, payload: CreateCommentPayload) => {
    try {
      const comment = await commentService.createComment(taskId, payload);
      setComments((prev) => [...prev, comment]);
      return comment;
    } catch (err) {
      toast.error(getErrorMessage(err, "Gửi bình luận thất bại"));
      throw err;
    }
  };

  const deleteComment = async (taskId: number, commentId: number) => {
    try {
      await commentService.deleteComment(taskId, commentId);
      setComments((prev) => prev.filter((c) => c.commentId !== commentId));
    } catch (err) {
      toast.error(getErrorMessage(err, "Xóa bình luận thất bại"));
      throw err;
    }
  };

  return { comments, loading, fetchComments, createComment, deleteComment };
}