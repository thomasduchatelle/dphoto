/**
 * A thunk is a function triggered by a user interaction which execute the business logic (API calls, actions, ...). It takes all its dependencies as arguments.
 *
 * Declaring the thunks with these interface is used to convert it into a stable handler ready to be called by the UI callbacks.
 */

export interface ThunkFactoryExtraArgs<PartialState> {
    partialState: PartialState
}

/**
 * ThunkDeclaration describes a thunk with a name, a factory, and a selector.
 * S: State interface
 * A: Action interface (dispatch signature)
 * F: The type of the callback function returned by the thunk factory
 * P: The type of the partial state returned by the selector
 */
export interface ThunkDeclaration<State, PartialState, F extends Function, FactoryArgs> {
    factory: (args: FactoryArgs & ThunkFactoryExtraArgs<PartialState>) => F;
    selector: (state: State) => PartialState;
}

type ActionCreator<TAction, TArgs extends any[]> = (...args: TArgs) => TAction;

/**
 * Shorthand to implement the declaration of a thunk dispatching immediately and always an action. The thunk tooling gives a stable handler with the dispatch injected.
 * @param actionCreator function to create the action (using createAction from daction lib)
 */
export function createSimpleThunkDeclaration<TAction, TArgs extends any[]>(
    actionCreator: ActionCreator<TAction, TArgs>
): ThunkDeclaration<any, {}, (...args: TArgs) => void, { dispatch: (action: TAction) => void }> {
    return {
        selector: () => ({}),
        factory: ({dispatch}) => {
            return (...args: TArgs) => {
                dispatch(actionCreator(...args));
            };
        },
    };
}