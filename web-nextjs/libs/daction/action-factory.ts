/**
 * Action is inspired by Redux actions, but with a reducer function included.
 *
 * Design goal is to extend reducer(s) without having to modify them: the function is included. It is aimed at reduce the context required by LLM to work.
 */

export interface Action<TState = any, TPayload = void> {
    type: string;
    reducer: (state: TState, action: Action<TState, TPayload>) => TState;
    payload?: TPayload;
}

// Overloaded function signatures for different payload types

export function createAction<TState>(
    type: string,
    reducer: (state: TState) => TState
): () => Action<TState, void>;

export function createAction<TState, TPayload>(
    type: string,
    reducer: (state: TState, payload: TPayload) => TState
): (payload: TPayload) => Action<TState, TPayload>;

export function createAction<TState, TPayload extends readonly any[]>(
    type: string,
    reducer: (state: TState, ...payload: TPayload) => TState
): (...payload: TPayload) => Action<TState, TPayload>;

// Implementation

export function createAction<TState, TPayload = void>(
    type: string,
    reducer: (state: TState, ...args: any[]) => TState
): any {
    const internalReducer = (state: TState, action: Action<TState, TPayload>) => {
        if (action.payload === undefined) {
            return reducer(state);
        } else if (Array.isArray(action.payload)) {
            return reducer(state, ...action.payload);
        } else {
            return reducer(state, action.payload);
        }
    };

    const actionCreator = (...args: any[]) => {
        const action: Action<TState, TPayload> = {
            type,
            reducer: internalReducer
        };

        if (args.length === 0) {
            // No payload
        } else if (args.length === 1) {
            action.payload = args[0] as TPayload;
        } else {
            action.payload = args as TPayload;
        }

        return action;
    };

    actionCreator.type = type;
    actionCreator.reducer = internalReducer;

    return actionCreator;
}

export function getPayload<TPayload>(action: Action<any, TPayload>): TPayload | undefined {
    return action.payload;
}