export interface ActionWithReducer<TState = any, TArgs extends any[] = any[]> {
    type: string;
    reducer: (state: TState, action: ActionWithReducer<TState, TArgs>) => TState;
    payload?: TArgs extends [] ? never : TArgs extends [infer T] ? T : TArgs;
}

export interface ActionCreator<TState, TArgs extends any[]> {
    (...args: TArgs): ActionWithReducer<TState, TArgs>;
    type: string;
    reducer: (state: TState, action: ActionWithReducer<TState, TArgs>) => TState;
}

export function createAction<TState, TArgs extends any[] = []>(
    type: string,
    reducer: (state: TState, action: ActionWithReducer<TState, TArgs>) => TState
): ActionCreator<TState, TArgs> {
    const actionCreator = ((...args: TArgs) => {
        const action: ActionWithReducer<TState, TArgs> = { 
            type,
            reducer
        };
        
        if (args.length === 1) {
            action.payload = args[0] as any;
        } else if (args.length > 1) {
            action.payload = args as any;
        }
        
        return action;
    }) as ActionCreator<TState, TArgs>;

    actionCreator.type = type;
    actionCreator.reducer = reducer;

    return actionCreator;
}

// Convenience functions for common patterns
export function createActionWithoutPayload<TState>(
    type: string,
    reducer: (state: TState, action: ActionWithReducer<TState, []>) => TState
): ActionCreator<TState, []> {
    return createAction(type, reducer);
}

export function createActionWithPayload<TState, TPayload>(
    type: string,
    reducer: (state: TState, action: ActionWithReducer<TState, [TPayload]>) => TState
): ActionCreator<TState, [TPayload]> {
    return createAction(type, reducer);
}

export function createActionWith2Payloads<TState, T1, T2>(
    type: string,
    reducer: (state: TState, action: ActionWithReducer<TState, [T1, T2]>) => TState
): ActionCreator<TState, [T1, T2]> {
    return createAction(type, reducer);
}

export function createActionWith3Payloads<TState, T1, T2, T3>(
    type: string,
    reducer: (state: TState, action: ActionWithReducer<TState, [T1, T2, T3]>) => TState
): ActionCreator<TState, [T1, T2, T3]> {
    return createAction(type, reducer);
}
