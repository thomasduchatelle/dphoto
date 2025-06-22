export interface ActionWithReducer<TState = any, TPayload = void> {
    type: string;
    reducer: (state: TState, action: ActionWithReducer<TState, TPayload>) => TState;
    payload?: TPayload;
}

// Overloaded function signatures for different payload types
export function createAction<TState>(
    type: string,
    reducer: (state: TState) => TState
): () => ActionWithReducer<TState, void>;

export function createAction<TState, TPayload>(
    type: string,
    reducer: (state: TState, payload: TPayload) => TState
): (payload: TPayload) => ActionWithReducer<TState, TPayload>;

export function createAction<TState, TPayload extends readonly any[]>(
    type: string,
    reducer: (state: TState, ...payload: TPayload) => TState
): (...payload: TPayload) => ActionWithReducer<TState, TPayload>;

// Implementation
export function createAction<TState, TPayload = void>(
    type: string,
    reducer: (state: TState, ...args: any[]) => TState
): any {
    const internalReducer = (state: TState, action: ActionWithReducer<TState, TPayload>) => {
        if (action.payload === undefined) {
            return reducer(state);
        } else if (Array.isArray(action.payload)) {
            return reducer(state, ...action.payload);
        } else {
            return reducer(state, action.payload);
        }
    };

    const actionCreator = (...args: any[]) => {
        const action: ActionWithReducer<TState, TPayload> = {
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

// Legacy convenience functions for backward compatibility
export function createActionWithoutPayload<TState>(
    type: string,
    reducer: (state: TState) => TState
): () => ActionWithReducer<TState, void> {
    return createAction(type, reducer);
}

export function createActionWithPayload<TState, TPayload>(
    type: string,
    reducer: (state: TState, payload: TPayload) => TState
): (payload: TPayload) => ActionWithReducer<TState, TPayload> {
    return createAction(type, reducer);
}
