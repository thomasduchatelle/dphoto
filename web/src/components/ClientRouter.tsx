'use client';

import {ReactNode, useEffect, useState} from 'react';

export interface RouterContextValue {
    path: string;
    params: Record<string, string>;
    query: URLSearchParams;
    navigate: (path: string) => void;
    replace: (path: string) => void;
}

export function useClientRouter(): RouterContextValue {
    // Use browser location instead of Waku router to enable SPA navigation
    const [currentPath, setCurrentPath] = useState(() => {
        if (typeof window !== 'undefined') {
            return window.location.pathname;
        }
        return '/';
    });

    useEffect(() => {
        // Listen for popstate events (browser back/forward)
        const handlePopState = () => {
            setCurrentPath(window.location.pathname);
        };

        window.addEventListener('popstate', handlePopState);
        return () => window.removeEventListener('popstate', handlePopState);
    }, []);

    const getCurrentPath = () => {
        return currentPath;
    };

    const getCurrentQuery = () => {
        if (typeof window !== 'undefined') {
            return new URLSearchParams(window.location.search);
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
        }
    };

    const replace = (newPath: string) => {
        if (typeof window !== 'undefined') {
            // Use replaceState to update URL without reload
            window.history.replaceState({}, '', newPath);
            setCurrentPath(window.location.pathname);
        }
    };

    return {
        path: getCurrentPath(),
        params: getCurrentParams(),
        query: getCurrentQuery(),
        navigate,
        replace,
    };
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
