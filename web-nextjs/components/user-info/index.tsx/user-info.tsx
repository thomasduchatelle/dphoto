import Image from 'next/image';
import Link from 'next/link';
import {basePath} from '@/libs/requests';

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
                href={`${basePath}/auth/logout`}
                className="ml-2 px-3 py-1 text-xs font-medium text-zinc-700 dark:text-zinc-300 hover:text-zinc-900 dark:hover:text-zinc-100 transition-colors"
            >
                Logout
            </Link>
        </div>
    );
}
