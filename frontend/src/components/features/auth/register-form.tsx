"use client";

import { useAuth } from "@/hooks/use-auth";
import { cn } from "@/lib/utils";
import { useRouter } from "next/navigation";
import { useState } from "react";

export function RegisterForm() {
  const router = useRouter();
  const { register, isLoading } = useAuth();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (password !== confirmPassword) {
      setError("パスワードが一致しません");
      return;
    }

    if (password.length < 8) {
      setError("パスワードは8文字以上で入力してください");
      return;
    }

    const result = await register({ email, password, name });
    if (result.success) {
      router.push("/dashboard");
    } else {
      setError(result.error?.error.message || "登録に失敗しました");
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4 w-full max-w-sm">
      <div>
        <label htmlFor="name" className="block text-sm font-medium mb-1">
          名前
        </label>
        <input
          id="name"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
          className={cn(
            "w-full px-4 py-2 border rounded-lg",
            "focus:outline-none focus:ring-2 focus:ring-blue-500",
            "dark:bg-gray-800 dark:border-gray-700"
          )}
          placeholder="山田 太郎"
        />
      </div>

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
          minLength={8}
          className={cn(
            "w-full px-4 py-2 border rounded-lg",
            "focus:outline-none focus:ring-2 focus:ring-blue-500",
            "dark:bg-gray-800 dark:border-gray-700"
          )}
          placeholder="8文字以上"
        />
      </div>

      <div>
        <label htmlFor="confirmPassword" className="block text-sm font-medium mb-1">
          パスワード（確認）
        </label>
        <input
          id="confirmPassword"
          type="password"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          required
          className={cn(
            "w-full px-4 py-2 border rounded-lg",
            "focus:outline-none focus:ring-2 focus:ring-blue-500",
            "dark:bg-gray-800 dark:border-gray-700"
          )}
          placeholder="パスワードを再入力"
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
        {isLoading ? "登録中..." : "新規登録"}
      </button>

      <p className="text-center text-sm text-gray-600 dark:text-gray-400">
        すでにアカウントをお持ちの方は{" "}
        <a href="/login" className="text-blue-600 hover:underline">
          ログイン
        </a>
      </p>
    </form>
  );
}
