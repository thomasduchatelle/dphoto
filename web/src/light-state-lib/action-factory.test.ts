import {createAction} from './action-factory';

interface DeathStarState {
    isOperational: boolean;
    shieldGenerator: {
        isActive: boolean;
        powerLevel: number;
    };
    reactorCore: {
        temperature: number;
        isStable: boolean;
    };
    hangarBays: Array<{
        id: string;
        occupiedBy?: string;
    }>;
    currentCommander?: string;
}

const initialDeathStarState: DeathStarState = {
    isOperational: false,
    shieldGenerator: {
        isActive: false,
        powerLevel: 0,
    },
    reactorCore: {
        temperature: 20,
        isStable: true,
    },
    hangarBays: [
        { id: 'bay-1' },
        { id: 'bay-2' },
        { id: 'bay-3' },
    ],
};

describe('action-factory', () => {
    describe('createAction with no payload', () => {
        const activateDeathStar = createAction<DeathStarState>(
            'ActivateDeathStar',
            (state: DeathStarState) => ({
                ...state,
                isOperational: true,
                shieldGenerator: {
                    ...state.shieldGenerator,
                    isActive: true,
                    powerLevel: 100,
                },
            })
        );

        it('creates action without payload', () => {
            const action = activateDeathStar();
            
            expect(action.type).toBe('ActivateDeathStar');
            expect(action.payload).toBeUndefined();
            expect(typeof action.reducer).toBe('function');
        });

        it('executes reducer correctly', () => {
            const action = activateDeathStar();
            const newState = action.reducer(initialDeathStarState, action);
            
            expect(newState.isOperational).toBe(true);
            expect(newState.shieldGenerator.isActive).toBe(true);
            expect(newState.shieldGenerator.powerLevel).toBe(100);
        });

        it('supports action comparison', () => {
            const action1 = activateDeathStar();
            const action2 = activateDeathStar();
            
            expect(action1).toEqual(action2);
            expect([action1]).toContainEqual(action2);
        });
    });

    describe('createAction with single payload', () => {
        const setReactorTemperature = createAction<DeathStarState, number>(
            'SetReactorTemperature',
            (state: DeathStarState, temperature: number) => ({
                ...state,
                reactorCore: {
                    ...state.reactorCore,
                    temperature,
                    isStable: temperature < 100,
                },
            })
        );

        it('creates action with payload', () => {
            const action = setReactorTemperature(85);
            
            expect(action.type).toBe('SetReactorTemperature');
            expect(action.payload).toBe(85);
            expect(typeof action.reducer).toBe('function');
        });

        it('executes reducer with payload', () => {
            const action = setReactorTemperature(150);
            const newState = action.reducer(initialDeathStarState, action);
            
            expect(newState.reactorCore.temperature).toBe(150);
            expect(newState.reactorCore.isStable).toBe(false);
        });

        it('distinguishes actions with different payloads', () => {
            const action1 = setReactorTemperature(50);
            const action2 = setReactorTemperature(75);
            
            expect(action1).not.toEqual(action2);
            expect([action1]).not.toContainEqual(action2);
        });

        it('compares actions with same payload', () => {
            const action1 = setReactorTemperature(50);
            const action2 = setReactorTemperature(50);
            
            expect(action1).toEqual(action2);
            expect([action1]).toContainEqual(action2);
        });
    });

    describe('createAction with tuple payload', () => {
        const assignShipToHangar = createAction<DeathStarState, [string, string]>(
            'AssignShipToHangar',
            (state: DeathStarState, hangarId: string, shipName: string) => ({
                ...state,
                hangarBays: state.hangarBays.map(bay =>
                    bay.id === hangarId
                        ? { ...bay, occupiedBy: shipName }
                        : bay
                ),
            })
        );

        it('creates action with tuple payload', () => {
            const action = assignShipToHangar('bay-1', 'TIE Fighter');
            
            expect(action.type).toBe('AssignShipToHangar');
            expect(action.payload).toEqual(['bay-1', 'TIE Fighter']);
            expect(typeof action.reducer).toBe('function');
        });

        it('executes reducer with multiple parameters', () => {
            const action = assignShipToHangar('bay-2', 'Imperial Shuttle');
            const newState = action.reducer(initialDeathStarState, action);
            
            const bay2 = newState.hangarBays.find(bay => bay.id === 'bay-2');
            expect(bay2?.occupiedBy).toBe('Imperial Shuttle');
            
            // Other bays should remain unchanged
            const bay1 = newState.hangarBays.find(bay => bay.id === 'bay-1');
            expect(bay1?.occupiedBy).toBeUndefined();
        });

        it('compares tuple actions correctly', () => {
            const action1 = assignShipToHangar('bay-1', 'TIE Fighter');
            const action2 = assignShipToHangar('bay-1', 'TIE Fighter');
            const action3 = assignShipToHangar('bay-2', 'TIE Fighter');
            
            expect(action1).toEqual(action2);
            expect(action1).not.toEqual(action3);
            expect([action1]).toContainEqual(action2);
            expect([action1]).not.toContainEqual(action3);
        });
    });

    describe('createAction with complex payload', () => {
        interface CommanderAssignment {
            name: string;
            rank: string;
            clearanceLevel: number;
        }

        const assignCommander = createAction<DeathStarState, CommanderAssignment>(
            'AssignCommander',
            (state: DeathStarState, commander: CommanderAssignment) => ({
                ...state,
                currentCommander: `${commander.rank} ${commander.name}`,
                shieldGenerator: {
                    ...state.shieldGenerator,
                    powerLevel: commander.clearanceLevel >= 9 ? 100 : 50,
                },
            })
        );

        it('handles complex object payloads', () => {
            const commander = {
                name: 'Tarkin',
                rank: 'Grand Moff',
                clearanceLevel: 10,
            };
            
            const action = assignCommander(commander);
            const newState = action.reducer(initialDeathStarState, action);
            
            expect(newState.currentCommander).toBe('Grand Moff Tarkin');
            expect(newState.shieldGenerator.powerLevel).toBe(100);
        });

        it('compares complex payload actions', () => {
            const commander1 = { name: 'Tarkin', rank: 'Grand Moff', clearanceLevel: 10 };
            const commander2 = { name: 'Tarkin', rank: 'Grand Moff', clearanceLevel: 10 };
            const commander3 = { name: 'Vader', rank: 'Lord', clearanceLevel: 10 };
            
            const action1 = assignCommander(commander1);
            const action2 = assignCommander(commander2);
            const action3 = assignCommander(commander3);
            
            expect(action1).toEqual(action2);
            expect(action1).not.toEqual(action3);
        });
    });

    describe('integration with generic reducer', () => {
        const activateDeathStar = createAction<DeathStarState>(
            'ActivateDeathStar',
            (state: DeathStarState) => ({ ...state, isOperational: true })
        );

        const setTemperature = createAction<DeathStarState, number>(
            'SetTemperature',
            (state: DeathStarState, temp: number) => ({
                ...state,
                reactorCore: { ...state.reactorCore, temperature: temp }
            })
        );

        // Simulate the generic reducer pattern
        function genericReducer(
            state: DeathStarState,
            action: any
        ): DeathStarState {
            if ('reducer' in action && typeof action.reducer === 'function') {
                return action.reducer(state, action);
            }
            return state;
        }

        it('works with generic reducer pattern', () => {
            let state = initialDeathStarState;
            
            // Apply activation
            const activateAction = activateDeathStar();
            state = genericReducer(state, activateAction);
            expect(state.isOperational).toBe(true);
            
            // Apply temperature change
            const tempAction = setTemperature(75);
            state = genericReducer(state, tempAction);
            expect(state.reactorCore.temperature).toBe(75);
        });
    });
});
