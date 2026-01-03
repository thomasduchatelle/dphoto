'use client';

import { useUser } from '@/lib/context/user-context';
import Image from 'next/image';

/**
 * Client component that displays the current user's information.
 * Uses the UserContext to access user data.
 */
export function UserDisplay() {
    const { user } = useUser();

    if (!user) {
        return (
            <div className="rounded-lg border border-zinc-200 bg-white p-4 shadow-sm dark:border-zinc-800 dark:bg-zinc-900">
                <p className="text-sm text-zinc-600 dark:text-zinc-400">
                    No user information available
                </p>
            </div>
        );
    }

    return (
        <div className="rounded-lg border border-zinc-200 bg-white p-6 shadow-sm dark:border-zinc-800 dark:bg-zinc-900">
            <h2 className="mb-4 text-lg font-semibold text-zinc-900 dark:text-zinc-50">
                User Information
            </h2>
            <div className="flex items-start gap-4">
                {user.picture && (
                    <Image
                        src={user.picture}
                        alt={user.name}
                        width={64}
                        height={64}
                        className="rounded-full"
                    />
                )}
                <div className="flex-1">
                    <div className="mb-2">
                        <span className="text-sm font-medium text-zinc-700 dark:text-zinc-300">
                            Name:
                        </span>
                        <span className="ml-2 text-sm text-zinc-900 dark:text-zinc-100">
                            {user.name}
                        </span>
                    </div>
                    <div className="mb-2">
                        <span className="text-sm font-medium text-zinc-700 dark:text-zinc-300">
                            Email:
                        </span>
                        <span className="ml-2 text-sm text-zinc-900 dark:text-zinc-100">
                            {user.email}
                        </span>
                    </div>
                    <div>
                        <span className="text-sm font-medium text-zinc-700 dark:text-zinc-300">
                            Role:
                        </span>
                        <span className="ml-2 text-sm text-zinc-900 dark:text-zinc-100">
                            {user.isOwner ? 'Owner' : 'User'}
                        </span>
                    </div>
                </div>
            </div>
        </div>
    );
}
