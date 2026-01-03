'use client';

import {useRouter} from "waku";

export interface RouterContextValue {
    path: string;
    params: Record<string, string>;
    query: URLSearchParams;
    navigate: (path: string) => void;
    replace: (path: string) => void;
}

export function useClientRouter(): RouterContextValue {
    const router = useRouter();
    return {
        navigate: (path: string) => router.push(path),
        params: parseParams(router.path),
        path: router.path,
        query: new URLSearchParams(router.query),
        replace: (path: string) => router.replace(path),
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
