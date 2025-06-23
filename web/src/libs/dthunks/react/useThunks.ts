import {useCallback, useMemo} from "react";
import {useStableSelector} from "./useStableSelector";
import {ThunkDeclaration} from "../api";

/**
 * useThunks aggregates thunks from a map of ThunkDeclaration.
 * Returns an object with the same keys, each mapped to its callback function.
 * The returned object is stable as long as all the callbacks are the same.
 */
export function useThunks<
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
    // Compute all callbacks
    const callbacks = {} as {
        [K in keyof T]: T[K] extends ThunkDeclaration<State, any, infer F, PartialFactoryArgs> ? F : never
    };
    for (const key in declarations) {
        const declaration = declarations[key];
        // eslint-disable-next-line react-hooks/rules-of-hooks
        const stableSelector = useStableSelector(declaration.selector);
        const partialState = stableSelector(state);
        // eslint-disable-next-line react-hooks/rules-of-hooks, react-hooks/exhaustive-deps
        callbacks[key] = useCallback(
            declaration.factory({
                ...factoryArgs,
                partialState
            }),
            [factoryArgs, partialState]
        ) as any;
    }
    // Memoize the returned object so its reference is stable if all callbacks are the same
    // eslint-disable-next-line react-hooks/rules-of-hooks, react-hooks/exhaustive-deps
    return useMemo(() => callbacks, Object.values(callbacks));
}
