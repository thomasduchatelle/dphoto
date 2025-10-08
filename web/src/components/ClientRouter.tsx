'use client';

import {ReactNode, useEffect, useState, createContext, useContext} from 'react';

export interface RouterContextValue {
    path: string;
    params: Record<string, string>;
    query: URLSearchParams;
    navigate: (path: string) => void;
    replace: (path: string) => void;
}

const RouterContext = createContext<RouterContextValue | null>(null);

export function RouterProvider({children}: {children: ReactNode}) {
    // Use browser location instead of Waku router to enable SPA navigation
    const [currentPath, setCurrentPath] = useState(() => {
        if (typeof window !== 'undefined') {
            return window.location.pathname;
        }
        return '/';
    });

    // Track query params separately to force re-render when they change
    const [currentSearch, setCurrentSearch] = useState(() => {
        if (typeof window !== 'undefined') {
            return window.location.search;
        }
        return '';
    });

    useEffect(() => {
        // Listen for popstate events (browser back/forward)
        const handlePopState = () => {
            setCurrentPath(window.location.pathname);
            setCurrentSearch(window.location.search);
        };

        window.addEventListener('popstate', handlePopState);
        return () => window.removeEventListener('popstate', handlePopState);
    }, []);

    const getCurrentPath = () => {
        return currentPath;
    };

    const getCurrentQuery = () => {
        if (typeof window !== 'undefined') {
            return new URLSearchParams(currentSearch);
        }
        return new URLSearchParams();
    };

    const getCurrentParams = () => {
        return parseParams(currentPath);
    };

    const navigate = (newPath: string) => {
        if (typeof window !== 'undefined') {
            // Use pushState to update URL without reload
            window.history.pushState({}, '', newPath);
            setCurrentPath(window.location.pathname);
            setCurrentSearch(window.location.search);
        }
    };

    const replace = (newPath: string) => {
        if (typeof window !== 'undefined') {
            // Use replaceState to update URL without reload
            window.history.replaceState({}, '', newPath);
            setCurrentPath(window.location.pathname);
            setCurrentSearch(window.location.search);
        }
    };

    const value: RouterContextValue = {
        path: getCurrentPath(),
        params: getCurrentParams(),
        query: getCurrentQuery(),
        navigate,
        replace,
    };

    return (
        <RouterContext.Provider value={value}>
            {children}
        </RouterContext.Provider>
    );
}

export function useClientRouter(): RouterContextValue {
    const context = useContext(RouterContext);
    if (!context) {
        throw new Error('useClientRouter must be used within a RouterProvider');
    }
    return context;
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

export function ClientLink({to, children, className, onClick}: {
    to: string;
    children: ReactNode;
    className?: string;
    onClick?: (e: React.MouseEvent) => void;
}) {
    const {navigate} = useClientRouter();

    const handleClick = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();
        if (onClick) {
            onClick(e);
        }
        navigate(to);
    };

    return (
        <a href={to} onClick={handleClick} className={className}>
            {children}
        </a>
    );
}
