import api from "./axios";
import type {
  RegisterPayload,
  LoginPayload,
  LoginResponse,
  RegisterResponse,
  RefreshTokenResponse,
} from "../types/user";
import type { User } from "../types/user";

export const authService = {
  register: async (payload: RegisterPayload): Promise<RegisterResponse> => {
    const res = await api.post<RegisterResponse>("/auth/register", payload);
    return res.data;
  },

  login: async (payload: LoginPayload): Promise<LoginResponse> => {
    const res = await api.post<LoginResponse>("/auth/login", payload);
    return res.data;
  },

  refresh: async (refreshToken: string): Promise<RefreshTokenResponse> => {
    const res = await api.post<RefreshTokenResponse>("/auth/refresh", {
      refreshToken,
    });
    return res.data;
  },

  logout: async (refreshToken: string): Promise<void> => {
    await api.post("/auth/logout", { refreshToken });
  },

  getMe: async (): Promise<User> => {
    const res = await api.get<User>("/auth/me");
    return res.data;
  },
};