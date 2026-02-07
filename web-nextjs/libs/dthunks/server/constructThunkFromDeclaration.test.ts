import {describe, expect, it} from 'vitest';
import {constructThunkFromDeclaration} from './constructThunkFromDeclaration';
import {ThunkDeclaration} from '../api';
import {Action} from '@/libs/daction';

interface TestState {
    counter: number;
    items: string[];
}

const incrementAction = (amount: number): Action<TestState, { amount: number }> => ({
    type: 'increment',
    payload: {amount},
    reducer: (state) => ({
        ...state,
        counter: state.counter + amount,
    })
});

const addItemAction = (item: string): Action<TestState, { item: string }> => ({
    type: 'addItem',
    payload: {item},
    reducer: (state) => ({
        ...state,
        items: [...state.items, item],
    })
});

describe('constructThunkFromDeclaration', () => {
    it('should execute thunk and return final state after reducing all dispatched actions', async () => {
        const thunkDeclaration: ThunkDeclaration<
            TestState,
            { counter: number },
            (value: number) => Promise<void>,
            { dispatch: (action: Action<TestState, any>) => void }
        > = {
            factory: ({dispatch, partialState}) => {
                return async (value: number) => {
                    dispatch(incrementAction(partialState.counter + value));
                    dispatch(addItemAction(`item-${value}`));
                };
            },
            selector: (state: TestState) => ({counter: state.counter})
        };

        const initialState: TestState = {counter: 5, items: []};
        const factoryArgs = {};

        const thunkExecutor = constructThunkFromDeclaration(
            thunkDeclaration,
            initialState,
            factoryArgs
        );

        const finalState = await thunkExecutor(10);

        expect(finalState.counter).toBe(20);
        expect(finalState.items).toEqual(['item-10']);
    });

    it('should handle multiple action dispatches', async () => {
        const thunkDeclaration: ThunkDeclaration<
            TestState,
            {},
            () => Promise<void>,
            { dispatch: (action: Action<TestState, any>) => void }
        > = {
            factory: ({dispatch}) => {
                return async () => {
                    dispatch(incrementAction(1));
                    dispatch(incrementAction(2));
                    dispatch(addItemAction('first'));
                    dispatch(incrementAction(3));
                    dispatch(addItemAction('second'));
                };
            },
            selector: () => ({})
        };

        const initialState: TestState = {counter: 0, items: []};
        const thunkExecutor = constructThunkFromDeclaration(
            thunkDeclaration,
            initialState,
            {}
        );

        const finalState = await thunkExecutor();

        expect(finalState.counter).toBe(6);
        expect(finalState.items).toEqual(['first', 'second']);
    });

    it('should propagate errors from thunk execution', async () => {
        const thunkDeclaration: ThunkDeclaration<
            TestState,
            {},
            () => Promise<void>,
            { dispatch: (action: Action<TestState, any>) => void }
        > = {
            factory: () => {
                return async () => {
                    throw new Error('Thunk execution failed');
                };
            },
            selector: () => ({})
        };

        const initialState: TestState = {counter: 0, items: []};
        const thunkExecutor = constructThunkFromDeclaration(
            thunkDeclaration,
            initialState,
            {}
        );

        await expect(thunkExecutor()).rejects.toThrow('Thunk execution failed');
    });

    it('should handle thunks with no dispatched actions', async () => {
        const thunkDeclaration: ThunkDeclaration<
            TestState,
            {},
            () => Promise<void>,
            { dispatch: (action: Action<TestState, any>) => void }
        > = {
            factory: () => {
                return async () => {
                };
            },
            selector: () => ({})
        };

        const initialState: TestState = {counter: 42, items: ['existing']};
        const thunkExecutor = constructThunkFromDeclaration(
            thunkDeclaration,
            initialState,
            {}
        );

        const finalState = await thunkExecutor();

        expect(finalState).toEqual(initialState);
    });
});
