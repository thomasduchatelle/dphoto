import Link from 'next/link';
import {clearAuthCookies} from '@/libs/security';
import {basePath} from '@/libs/requests';

export default async function LogoutCallbackPage() {
    await clearAuthCookies();

    return (
        <div className="flex min-h-screen items-center justify-center bg-zinc-50 dark:bg-black">
            <div className="w-full max-w-md rounded-lg bg-white p-8 shadow-lg dark:bg-zinc-900 dark:shadow-zinc-800">
                <div className="mb-6 text-center">
                    <svg
                        className="mx-auto mb-4 h-16 w-16 text-green-500"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                    </svg>
                    <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-100">
                        Logged Out Successfully
                    </h1>
                    <p className="mt-2 text-zinc-600 dark:text-zinc-400">
                        You have been successfully logged out of your account.
                    </p>
                </div>
                <div className="mt-6">
                    <Link
                        href={`${basePath}/auth/login`}
                        className="flex w-full items-center justify-center rounded-full bg-zinc-900 px-6 py-3 text-white transition-colors hover:bg-zinc-700 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-300"
                    >
                        Log Back In
                    </Link>
                </div>
            </div>
        </div>
    );
}
