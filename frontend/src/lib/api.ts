import type { ApiError, AuthResponse, LoginRequest, RegisterRequest, User } from "@/types/auth";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

class ApiClient {
  private accessToken: string | null = null;

  setAccessToken(token: string | null) {
    this.accessToken = token;
  }

  getAccessToken() {
    return this.accessToken;
  }

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;

    const headers: HeadersInit = {
      "Content-Type": "application/json",
      ...options.headers,
    };

    if (this.accessToken) {
      (headers as Record<string, string>).Authorization = `Bearer ${this.accessToken}`;
    }

    const response = await fetch(url, {
      ...options,
      headers,
      credentials: "include", // For cookies (refresh token)
    });

    const data = await response.json();

    if (!response.ok) {
      throw data as ApiError;
    }

    return data as T;
  }

  // Auth endpoints
  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await this.request<AuthResponse>("/auth/register", {
      method: "POST",
      body: JSON.stringify(data),
    });
    this.setAccessToken(response.data.access_token);
    return response;
  }

  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await this.request<AuthResponse>("/auth/login", {
      method: "POST",
      body: JSON.stringify(data),
    });
    this.setAccessToken(response.data.access_token);
    return response;
  }

  async refresh(): Promise<AuthResponse> {
    const response = await this.request<AuthResponse>("/auth/refresh", {
      method: "POST",
    });
    this.setAccessToken(response.data.access_token);
    return response;
  }

  async logout(): Promise<void> {
    await this.request("/auth/logout", {
      method: "POST",
    });
    this.setAccessToken(null);
  }

  async getMe(): Promise<{ success: boolean; data: User }> {
    return this.request("/auth/me");
  }
}

export const apiClient = new ApiClient();
