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

export const projectSchema = z.object({
  name: z.string().min(1, "Tên project là bắt buộc").max(200),
  description: z.string().max(1000).optional(),
  deadline: z.string().optional(),
});
export type ProjectFormData = z.infer<typeof projectSchema>;

export const epicSchema = z.object({
  title: z.string().min(1, "Tiêu đề là bắt buộc").max(200),
  description: z.string().max(1000).optional(),
});
export type EpicFormData = z.infer<typeof epicSchema>;

export const milestoneSchema = z.object({
  title: z.string().min(1, "Tiêu đề là bắt buộc").max(200),
  description: z.string().max(1000).optional(),
  dueDate: z.string().optional(),
});
export type MilestoneFormData = z.infer<typeof milestoneSchema>;

export const addMemberSchema = z.object({
  email: z.string().min(1, "Email là bắt buộc").email("Email không hợp lệ"),
});
export type AddMemberFormData = z.infer<typeof addMemberSchema>;

export const labelSchema = z.object({
  name: z.string().min(1, "Tên label là bắt buộc").max(50),
  color: z
    .string()
    .regex(/^#[0-9A-Fa-f]{6}$/, "Màu phải đúng định dạng hex, ví dụ #ff0000")
    .optional(),
});
export type LabelFormData = z.infer<typeof labelSchema>;

export const sprintSchema = z.object({
  name: z.string().min(1, "Tên sprint là bắt buộc").max(200),
  startDate: z.string().optional(),
  endDate: z.string().optional(),
});
export type SprintFormData = z.infer<typeof sprintSchema>;