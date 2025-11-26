import { RegisterForm } from "@/components/features/auth/register-form";

export default function RegisterPage() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-8">
      <div className="text-center mb-8">
        <h1 className="text-3xl font-bold mb-2">新規登録</h1>
        <p className="text-gray-600 dark:text-gray-400">アカウントを作成してください</p>
      </div>
      <RegisterForm />
    </main>
  );
}
