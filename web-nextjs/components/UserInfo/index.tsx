import Image from 'next/image';
import Link from 'next/link';

export interface UserInfoProps {
    name: string;
    email: string;
    picture?: string;
}

export function UserInfo({ name, email, picture }: UserInfoProps) {
    return (
        <div className="fixed top-4 right-4 flex items-center gap-3 bg-white dark:bg-zinc-900 rounded-full shadow-lg px-4 py-2 border border-zinc-200 dark:border-zinc-800">
            {picture ? (
                <Image
                    src={picture}
                    alt={name}
                    width={32}
                    height={32}
                    className="rounded-full"
                />
            ) : (
                <div className="w-8 h-8 rounded-full bg-zinc-300 dark:bg-zinc-700 flex items-center justify-center">
                    <span className="text-sm font-semibold text-zinc-700 dark:text-zinc-300">
                        {name.charAt(0).toUpperCase()}
                    </span>
                </div>
            )}
            <div className="flex flex-col">
                <span className="text-sm font-semibold text-zinc-900 dark:text-zinc-100">
                    {name}
                </span>
                <span className="text-xs text-zinc-600 dark:text-zinc-400">
                    {email}
                </span>
            </div>
            <Link
                href="/auth/logout"
                className="ml-2 rounded-full p-2 text-zinc-600 hover:bg-zinc-100 dark:text-zinc-400 dark:hover:bg-zinc-800 transition-colors"
                title="Logout"
            >
                <svg
                    className="h-5 w-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                >
                    <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
                    />
                </svg>
            </Link>
        </div>
    );
}
