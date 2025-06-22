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