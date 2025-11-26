"use client";

import { useAuth } from "@/hooks/use-auth";
import { cn } from "@/lib/utils";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function DashboardPage() {
  const router = useRouter();
  const { user, isLoading, isAuthenticated, logout } = useAuth();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push("/login");
    }
  }, [isLoading, isAuthenticated, router]);

  const handleLogout = async () => {
    await logout();
    router.push("/");
  };

  if (isLoading) {
    return (
      <main className="flex min-h-screen items-center justify-center">
        <p>読み込み中...</p>
      </main>
    );
  }

  if (!isAuthenticated || !user) {
    return null;
  }

  return (
    <main className="min-h-screen p-8">
      <div className="max-w-4xl mx-auto">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold">ダッシュボード</h1>
          <button
            onClick={handleLogout}
            className={cn(
              "px-4 py-2 rounded-lg",
              "border border-gray-300 hover:bg-gray-100",
              "dark:border-gray-700 dark:hover:bg-gray-800",
              "transition-colors"
            )}
          >
            ログアウト
          </button>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold mb-4">ユーザー情報</h2>
          <dl className="space-y-2">
            <div className="flex">
              <dt className="w-24 text-gray-500">名前:</dt>
              <dd>{user.name}</dd>
            </div>
            <div className="flex">
              <dt className="w-24 text-gray-500">メール:</dt>
              <dd>{user.email}</dd>
            </div>
            <div className="flex">
              <dt className="w-24 text-gray-500">ID:</dt>
              <dd className="text-sm text-gray-400">{user.id}</dd>
            </div>
          </dl>
        </div>
      </div>
    </main>
  );
}
