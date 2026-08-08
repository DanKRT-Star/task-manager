import { z } from "zod";

export const loginSchema = z.object({
  email: z.string().min(1, "Email là bắt buộc").email("Email không hợp lệ"),
  password: z.string().min(1, "Mật khẩu là bắt buộc"),
});

export type LoginFormData = z.infer<typeof loginSchema>;

export const registerSchema = z.object({
  userName: z.string().min(2, "Tên phải có ít nhất 2 ký tự"),
  email: z.string().min(1, "Email là bắt buộc").email("Email không hợp lệ"),
  password: z.string().min(8, "Mật khẩu phải có ít nhất 8 ký tự"),
});

export type RegisterFormData = z.infer<typeof registerSchema>;


export const taskSchema = z.object({
  title: z.string().min(1, "Tiêu đề là bắt buộc").max(200, "Tối đa 200 ký tự"),
  description: z.string().max(1000, "Tối đa 1000 ký tự").optional(),
  deadline: z.string().min(1, "Deadline là bắt buộc"),
  status: z.enum(["pending", "in_progress", "done"]),
});

export type TaskFormData = z.infer<typeof taskSchema>;