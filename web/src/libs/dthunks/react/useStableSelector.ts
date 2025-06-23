import {useCallback, useRef} from "react";

/**
 * useStableSelector wrap a selector function to ensure the same reference is returned when the result is the same. It expects the values of each property to be equals (reference).
 */
export function useStableSelector<T, R>(selector: (input: T) => R): (input: T) => R {
    const lastResultRef = useRef<R | null>(null);

    return useCallback((input: T) => {
        const result = selector(input);
        if (lastResultRef.current === null || !shallowEqual(lastResultRef.current, result)) {
            lastResultRef.current = result;
            return result;
        }
        return lastResultRef.current;
    }, [lastResultRef, selector])
}

function shallowEqual(objA: any, objB: any): boolean {
    if (objA === objB) return true;
    if (typeof objA !== "object" || objA === null ||
        typeof objB !== "object" || objB === null) {
        return false;
    }
    const keysA = Object.keys(objA);
    const keysB = Object.keys(objB);
    if (keysA.length !== keysB.length) return false;
    for (let key of keysA) {
        if (objA[key] !== objB[key]) return false;
    }
    return true;
}

