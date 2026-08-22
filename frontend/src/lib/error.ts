import axios from "axios";

export function getErrorMessage(err: unknown, fallback: string): string {
  if (axios.isAxiosError(err)) {
    const data = err.response?.data as { error?: string } | undefined;
    if (data?.error) return data.error;
    if (!err.response) return "Network error — please check your connection or try again.";
    return fallback;
  }
  if (err instanceof Error) return err.message;
  return fallback;
}