import { createContext, useState, type ReactNode } from "react";
import { authService } from "../services/auth_service";
import type { LoginPayload, RegisterPayload, User } from "../types/user";

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (payload: LoginPayload) => Promise<void>;
  register: (payload: RegisterPayload) => Promise<void>;
  logout: () => Promise<void>;
}

// eslint-disable-next-line react-refresh/only-export-components
export const AuthContext = createContext<AuthContextType | undefined>(
  undefined
);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(() => {
    const saved = localStorage.getItem("user");
    return saved ? JSON.parse(saved) : null;
  });
  const [accessToken, setAccessToken] = useState<string | null>(() =>
    localStorage.getItem("accessToken")
  );
  const [isLoading, setIsLoading] = useState(false);

  const login = async (payload: LoginPayload) => {
    setIsLoading(true);
    try {
      const res = await authService.login(payload);
      localStorage.setItem("accessToken", res.accessToken);
      localStorage.setItem("refreshToken", res.refreshToken);
      setAccessToken(res.accessToken);

      // Backend /auth/login không trả user, gọi thêm /auth/me để lấy thông tin hiển thị
      const me = await authService.getMe();
      localStorage.setItem("user", JSON.stringify(me));
      setUser(me);
    } finally {
      setIsLoading(false);
    }
  };

  const register = async (payload: RegisterPayload) => {
    setIsLoading(true);
    try {
      const res = await authService.register(payload);
      // Đăng ký xong không tự đăng nhập — giữ đúng flow cũ: yêu cầu người dùng đăng nhập lại
      return res.user as unknown as void;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = async () => {
    const refreshToken = localStorage.getItem("refreshToken");
    try {
      if (refreshToken) {
        await authService.logout(refreshToken);
      }
    } catch {
      // Dù API logout lỗi (token đã hết hạn sẵn), vẫn dọn sạch state phía client
    } finally {
      localStorage.removeItem("accessToken");
      localStorage.removeItem("refreshToken");
      localStorage.removeItem("user");
      setAccessToken(null);
      setUser(null);
    }
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated: !!accessToken,
        isLoading,
        login,
        register,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}