import { LoginForm } from "@/components/features/auth/login-form";

export default function LoginPage() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-8">
      <div className="text-center mb-8">
        <h1 className="text-3xl font-bold mb-2">ログイン</h1>
        <p className="text-gray-600 dark:text-gray-400">アカウントにログインしてください</p>
      </div>
      <LoginForm />
    </main>
  );
}
