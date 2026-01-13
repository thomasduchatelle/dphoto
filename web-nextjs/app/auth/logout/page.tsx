import Link from 'next/link';

export default async function LogoutPage() {

    return (
        <div className="flex min-h-screen items-center justify-center bg-zinc-50 dark:bg-black">
            <div className="mx-4 max-w-md rounded-lg bg-white p-8 shadow-lg dark:bg-zinc-900 dark:border dark:border-zinc-800">
                <div className="mb-6 flex justify-center">
                    <svg
                        className="h-16 w-16 text-green-500"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                    >
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                    </svg>
                </div>
                <h1 className="mb-4 text-center text-2xl font-bold text-zinc-900 dark:text-zinc-50">
                    Successfully Logged Out
                </h1>
                <p className="mb-8 text-center text-zinc-600 dark:text-zinc-400">
                    You have been successfully logged out of your account.
                </p>
                <div className="flex justify-center">
                    <Link
                        href="/auth/login"
                        className="rounded-full bg-zinc-900 px-6 py-3 text-sm font-medium text-white transition-colors hover:bg-zinc-700 dark:bg-zinc-50 dark:text-zinc-900 dark:hover:bg-zinc-200"
                    >
                        Sign In Again
                    </Link>
                </div>
            </div>
        </div>
    );
}
