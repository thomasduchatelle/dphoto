import {Action} from "@/libs/daction";

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