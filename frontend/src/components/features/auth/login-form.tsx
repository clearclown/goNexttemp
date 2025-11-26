"use client";

import { useAuth } from "@/hooks/use-auth";
import { cn } from "@/lib/utils";
import { useRouter } from "next/navigation";
import { useState } from "react";

export function LoginForm() {
  const router = useRouter();
  const { login, isLoading } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    const result = await login({ email, password });
    if (result.success) {
      router.push("/dashboard");
    } else {
      setError(result.error?.error.message || "ログインに失敗しました");
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4 w-full max-w-sm">
      <div>
        <label htmlFor="email" className="block text-sm font-medium mb-1">
          メールアドレス
        </label>
        <input
          id="email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          className={cn(
            "w-full px-4 py-2 border rounded-lg",
            "focus:outline-none focus:ring-2 focus:ring-blue-500",
            "dark:bg-gray-800 dark:border-gray-700"
          )}
          placeholder="example@email.com"
        />
      </div>

      <div>
        <label htmlFor="password" className="block text-sm font-medium mb-1">
          パスワード
        </label>
        <input
          id="password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          className={cn(
            "w-full px-4 py-2 border rounded-lg",
            "focus:outline-none focus:ring-2 focus:ring-blue-500",
            "dark:bg-gray-800 dark:border-gray-700"
          )}
          placeholder="••••••••"
        />
      </div>

      {error && <p className="text-red-500 text-sm">{error}</p>}

      <button
        type="submit"
        disabled={isLoading}
        className={cn(
          "w-full py-2 px-4 rounded-lg font-medium",
          "bg-blue-600 text-white hover:bg-blue-700",
          "disabled:opacity-50 disabled:cursor-not-allowed",
          "transition-colors"
        )}
      >
        {isLoading ? "ログイン中..." : "ログイン"}
      </button>

      <p className="text-center text-sm text-gray-600 dark:text-gray-400">
        アカウントをお持ちでない方は{" "}
        <a href="/register" className="text-blue-600 hover:underline">
          新規登録
        </a>
      </p>
    </form>
  );
}
