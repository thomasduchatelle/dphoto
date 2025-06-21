export interface ActionCreator<TPayload = void> {
    (...args: TPayload extends void ? [] : [TPayload]): Action<TPayload>;
    type: string;
    reducer: (state: any, ...args: TPayload extends void ? [any] : [any, TPayload]) => any;
}

export interface Action<TPayload = void> {
    type: string;
    payload?: TPayload;
}

export function createAction<TState, TPayload = void>(
    type: string,
    reducer: TPayload extends void 
        ? (state: TState, action: Action<void>) => TState
        : (state: TState, action: Action<TPayload>) => TState
): ActionCreator<TPayload> {
    const actionCreator = ((...args: any[]) => {
        const action: Action<TPayload> = { type };
        if (args.length > 0) {
            action.payload = args[0] as TPayload;
        }
        return action;
    }) as ActionCreator<TPayload>;

    actionCreator.type = type;
    actionCreator.reducer = reducer;

    return actionCreator;
}

export function createActionWithPayload<TState, TPayload>(
    type: string,
    reducer: (state: TState, action: Action<TPayload>) => TState
): ActionCreator<TPayload> {
    return createAction(type, reducer);
}

export function createActionWithoutPayload<TState>(
    type: string,
    reducer: (state: TState, action: Action<void>) => TState
): ActionCreator<void> {
    return createAction(type, reducer);
}
