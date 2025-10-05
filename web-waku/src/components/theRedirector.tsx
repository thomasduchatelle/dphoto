'use client';

import {useRouter} from "waku";
import {useEffect} from "react";

export function TheRedirector({count}: { count: number }) {
    const router = useRouter();

    useEffect(() => {
        router.replace(`/?counter=${count}`).then();
    }, [count]);
    return null;
}