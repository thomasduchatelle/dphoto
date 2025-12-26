'use client';

import {ReactNode, createContext, useContext, useState, useCallback} from 'react';
import {useRouter as useWakuRouter} from "waku";

export interface RouterContextValue {
    path: string;
    params: Record<string, string>;
    query: URLSearchParams;
    navigate: (path: string) => void;
    replace: (path: string) => void;
}

// Create a context for story/test environments
const RouterContext = createContext<RouterContextValue | null>(null);

export function RouterProvider({children}: { children: ReactNode }) {
    // For story environments where there's no Waku router
    const [mockPath, setMockPath] = useState('/');
    
    const navigate = useCallback((path: string) => {
        setMockPath(path);
        if (typeof window !== 'undefined' && window.history) {
            try {
                window.history.pushState({}, '', path);
            } catch (e) {
                // Ignore in test environments
            }
        }
    }, []);
    
    const replace = useCallback((path: string) => {
        setMockPath(path);
        if (typeof window !== 'undefined' && window.history) {
            try {
                window.history.replaceState({}, '', path);
            } catch (e) {
                // Ignore in test environments
            }
        }
    }, []);
    
    const mockRouter: RouterContextValue = {
        path: mockPath.split('?')[0],
        params: parseParams(mockPath),
        query: new URLSearchParams(mockPath.split('?')[1] || ''),
        navigate,
        replace,
    };
    
    return (
        <RouterContext.Provider value={mockRouter}>
            {children}
        </RouterContext.Provider>
    );
}

export function useClientRouter(): RouterContextValue {
    // First check if we have a story/test context
    const mockContext = useContext(RouterContext);
    if (mockContext) {
        return mockContext;
    }
    
    // Otherwise use the real Waku router
    try {
        const router = useWakuRouter();
        return {
            navigate: (path: string) => router.push(path),
            params: parseParams(router.path),
            path: router.path,
            query: new URLSearchParams(router.query),
            replace: (path: string) => router.replace(path),
        };
    } catch (error) {
        // If waku router is not available (e.g., in Ladle stories),
        // return a minimal fallback router
        return {
            path: '/',
            params: {},
            query: new URLSearchParams(),
            navigate: () => {},
            replace: () => {},
        };
    }
}

function parseParams(path: string): Record<string, string> {
    const cleanPath = path.split('?')[0];
    const parts = cleanPath.split('/').filter(p => p);

    // Match /albums/:owner/:album/:encodedId/:filename
    if (parts[0] === 'albums') {
        if (parts.length >= 5) {
            return {
                owner: parts[1],
                album: parts[2],
                encodedId: parts[3],
                filename: parts[4],
            };
        } else if (parts.length >= 3) {
            return {
                owner: parts[1],
                album: parts[2],
            };
        }
    }

    return {};
}
