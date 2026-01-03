'use client';

import React, { createContext, useContext, ReactNode } from 'react';
import { AuthenticatedUser } from '@/lib/security/constants';

interface UserContextValue {
    user: AuthenticatedUser | null;
}

const UserContext = createContext<UserContextValue | undefined>(undefined);

interface UserProviderProps {
    user: AuthenticatedUser | null;
    children: ReactNode;
}

/**
 * Client-side provider that makes user information available throughout the app.
 */
export function UserProvider({ user, children }: UserProviderProps) {
    return (
        <UserContext.Provider value={{ user }}>
            {children}
        </UserContext.Provider>
    );
}

/**
 * Hook to access the current user from any component.
 */
export function useUser(): UserContextValue {
    const context = useContext(UserContext);
    if (context === undefined) {
        throw new Error('useUser must be used within a UserProvider');
    }
    return context;
}
