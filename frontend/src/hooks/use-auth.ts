"use client";

import { apiClient } from "@/lib/api";
import type { ApiError, LoginRequest, RegisterRequest, User } from "@/types/auth";
import { useCallback, useEffect, useState } from "react";

interface AuthState {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
}

export function useAuth() {
  const [state, setState] = useState<AuthState>({
    user: null,
    isLoading: true,
    isAuthenticated: false,
  });

  const checkAuth = useCallback(async () => {
    try {
      // Try to refresh token on initial load
      const response = await apiClient.refresh();
      setState({
        user: response.data.user,
        isLoading: false,
        isAuthenticated: true,
      });
    } catch {
      setState({
        user: null,
        isLoading: false,
        isAuthenticated: false,
      });
    }
  }, []);

  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  const login = async (data: LoginRequest) => {
    setState((prev) => ({ ...prev, isLoading: true }));
    try {
      const response = await apiClient.login(data);
      setState({
        user: response.data.user,
        isLoading: false,
        isAuthenticated: true,
      });
      return { success: true };
    } catch (error) {
      setState((prev) => ({ ...prev, isLoading: false }));
      return { success: false, error: error as ApiError };
    }
  };

  const register = async (data: RegisterRequest) => {
    setState((prev) => ({ ...prev, isLoading: true }));
    try {
      const response = await apiClient.register(data);
      setState({
        user: response.data.user,
        isLoading: false,
        isAuthenticated: true,
      });
      return { success: true };
    } catch (error) {
      setState((prev) => ({ ...prev, isLoading: false }));
      return { success: false, error: error as ApiError };
    }
  };

  const logout = async () => {
    try {
      await apiClient.logout();
    } finally {
      setState({
        user: null,
        isLoading: false,
        isAuthenticated: false,
      });
    }
  };

  return {
    ...state,
    login,
    register,
    logout,
    checkAuth,
  };
}
