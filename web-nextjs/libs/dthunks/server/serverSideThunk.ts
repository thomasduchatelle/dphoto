import {ThunkDeclaration} from '../api';

import {ServerState} from "@/libs/dthunks/server/serverState";

/**
 * Creates a server-side thunk executor from a thunk declaration.
 * This function returns a callable that takes the initial state and thunk arguments,
 * executes the thunk, and returns the final computed state.
 *
 * @param declaration - The thunk declaration containing factory and selector
 * @param factoryArgs - Function that creates the adapter instance
 * @returns A function that takes initial state and thunk args, and returns the final state
 *
 * @example
 * const onPageRefresh = serverSideThunk(catalogThunks.onPageRefresh, newServerAdapterFactory)
 * const catalogState = await onPageRefresh(initialCatalogState(currentUser), undefined);
 */
export function serverSideThunk<TState, TPartialState, TArgs extends any[], TFactoryArgs extends { dispatch: any }>(
    declaration: ThunkDeclaration<TState, TPartialState, (...args: TArgs) => Promise<void>, TFactoryArgs>,
    factoryArgs: Omit<TFactoryArgs, "dispatch">,
) {
    return async (initialState: TState, ...args: TArgs): Promise<TState> => {
        const serverState = new ServerState(initialState);
        const partialState = declaration.selector(serverState.currentState());

        const thunkFunction = declaration.factory({
            ...factoryArgs,
            dispatch: serverState.dispatch,
            partialState,
        } as TFactoryArgs & { partialState: TPartialState });

        await thunkFunction(...args);
        return serverState.currentState();
    };
}
