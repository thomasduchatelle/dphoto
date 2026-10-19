"use client";

/* eslint-disable @typescript-eslint/no-explicit-any */

import {useMemo} from "react";
import {useRouter} from "next/navigation";
import {useThunks} from "@/libs/dthunks/react";
import {ThunkDeclaration} from "@/libs/dthunks";
import {isSessionExpiredException} from "./session-expired-exception";

export const SESSION_EXPIRED_PATH = "/auth/session-expired";

function wrapCallback<F>(fn: F, onSessionExpired: () => void): F {
    const wrapped = ((...args: any[]) => {
        let result: any;
        try {
            result = (fn as any)(...args);
        } catch (err) {
            if (isSessionExpiredException(err)) {
                onSessionExpired();
                return undefined;
            }
            throw err;
        }
        if (result && typeof result.then === "function") {
            return (result as Promise<any>).catch((err) => {
                if (isSessionExpiredException(err)) {
                    onSessionExpired();
                    return undefined;
                }
                throw err;
            });
        }
        return result;
    }) as unknown as F;
    return wrapped;
}

export function useSessionAwareThunks<
    State,
    PartialFactoryArgs,
    T extends { [K in keyof T]: ThunkDeclaration<State, any, any, PartialFactoryArgs> }
>(
    declarations: T,
    factoryArgs: PartialFactoryArgs,
    state: State,
): {
    [K in keyof T]: T[K] extends ThunkDeclaration<State, any, infer F, PartialFactoryArgs> ? F : never
} {
    const router = useRouter();
    const callbacks = useThunks(declarations, factoryArgs, state);

    return useMemo(() => {
        const onSessionExpired = () => {
            router.push(SESSION_EXPIRED_PATH);
        };
        const wrapped = {} as typeof callbacks;
        for (const key in callbacks) {
            wrapped[key] = wrapCallback(callbacks[key], onSessionExpired);
        }
        return wrapped;
    }, [callbacks, router]);
}
