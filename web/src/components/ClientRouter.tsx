'use client';

import {useRouter} from 'waku';
import {ReactNode} from 'react';

export interface RouterContextValue {
    path: string;
    params: Record<string, string>;
    query: URLSearchParams;
    navigate: (path: string) => void;
    replace: (path: string) => void;
}

export function useClientRouter(): RouterContextValue {
    const router = useRouter();

    const getCurrentPath = () => {
        return router.path;
    };

    const getCurrentQuery = () => {
        return new URLSearchParams(router.query);
    };

    const getCurrentParams = () => {
        return parseParams(router.path);
    };

    const navigate = (newPath: string) => {
        router.push(newPath);
    };

    const replace = (newPath: string) => {
        router.replace(newPath);
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
