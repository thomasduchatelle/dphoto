import {ThunkDeclaration} from '../api';
import {Action} from '@/libs/daction';

/**
 * Constructs an executable thunk from a declaration for server-side execution.
 * Unlike client-side thunks (useThunks hook), this executes synchronously and accumulates
 * all dispatched actions, then reduces them against the initial state to compute the final state.
 *
 * @param declaration - The thunk declaration containing factory and selector
 * @param factoryArgs - Factory arguments (adapters, etc.) provided to the thunk
 * @param serverState - The server state instance to manage dispatched actions and compute final state
 * @returns A function that executes the thunk and returns the final computed state
 */
export function constructThunkFromDeclaration<TState, TPartialState, TArgs extends any[], TFactoryArgs>(
    declaration: ThunkDeclaration<TState, TPartialState, (...args: TArgs) => Promise<void>, TFactoryArgs>,
    factoryArgs: TFactoryArgs,
    serverState: ServerState<TState, Action<TState, any>>,
): (...args: TArgs) => Promise<TState> {
    const partialState = declaration.selector(serverState.currentState());

    const thunkFunction = declaration.factory({
        ...factoryArgs,
        dispatch: serverState.dispatch,
        partialState,
    } as TFactoryArgs & { dispatch: typeof serverState.dispatch, partialState: TPartialState });

    return async (...args: TArgs): Promise<TState> => {
        await thunkFunction(...args);
        return serverState.currentState();
    };
}

export class ServerState<TState, A extends Action<TState, any>> {
    constructor(
        private readonly initialState: TState,
        private accumulatedActions: Action<TState, any>[] = [],
    ) {
    }

    public dispatch = (action: Action<TState, any>): void => {
        this.accumulatedActions.push(action);
    }

    public currentState = (): TState => {
        let finalState = this.initialState;
        for (const action of this.accumulatedActions) {
            finalState = action.reducer(finalState, action);
        }
        return finalState;
    }
}