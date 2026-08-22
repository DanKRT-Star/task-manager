export interface User {
  userId: number;
  userName: string;
  email: string;
  createdAt: string;
  updatedAt: string;
}

export interface RegisterPayload {
  userName: string;
  email: string;
  password: string;
}

export interface LoginPayload {
  email: string;
  password: string;
}

export interface LoginResponse {
  message: string;
  accessToken: string;
  refreshToken: string;
}

export interface RefreshTokenResponse {
  accessToken: string;
  refreshToken: string;
}

export interface RegisterResponse {
  message: string;
  user: User;
}