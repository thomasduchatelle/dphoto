import {Action} from "./action-factory";

export function createGenericReducer<TState>(): (state: TState, action: Action<TState>) => TState {
    return (state: TState, action: Action<TState> | { type: string }): TState => {
        if ('reducer' in action && typeof action.reducer === 'function') {
            return action.reducer(state, action as Action<TState>);
        }

        throw new Error(`Action does not have a reducer: ${action.type} [${JSON.stringify(action)}]`);
    };
}