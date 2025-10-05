import {renderHook} from "@testing-library/react";
import {useThunks} from "./useThunks";
import {ThunkDeclaration} from "../api";

// Dummy Star Wars state and types
type StarWarsState = {
    jedi: { name: string; lightsaber: string }[];
    sith: { name: string; lightsaber: string }[];
};

type ForceArg = { force: string };

// Dummy thunk factory and selector
function createJediGreetingFactory({partialState}: { partialState: { name: string } }) {
    return (greeting: string = "Hello") => `${greeting} ${partialState.name}, may the force be with you!`;
}

function firstJediSelector(state: StarWarsState) {
    return state.jedi[0];
}

function createSithGreetingFactory({partialState, force}: { partialState: { name: string }, force: string }) {
    return () => `Greetings ${partialState.name}, embrace the ${force}!`;
}

function firstSithSelector(state: StarWarsState) {
    return state.sith[0];
}

describe("useThunks", () => {
    const initialState: StarWarsState = {
        jedi: [{name: "Luke Skywalker", lightsaber: "green"}],
        sith: [{name: "Darth Vader", lightsaber: "red"}]
    };
    const declarations = {
        jediGreeting: {
            factory: createJediGreetingFactory,
            selector: firstJediSelector
        } as ThunkDeclaration<StarWarsState, { name: string }, (greeting?: string) => string, ForceArg>,
        sithGreeting: {
            factory: createSithGreetingFactory,
            selector: firstSithSelector
        } as ThunkDeclaration<StarWarsState, { name: string }, () => string, ForceArg>
    };

    it("should return an object with the same keys as the declarations when called with valid declarations", () => {
        const {result} = renderHook(() =>
            useThunks(declarations, {force: "Force"}, initialState)
        );
        expect(Object.keys(result.current)).toEqual(Object.keys(declarations));

        expect(result.current.jediGreeting("Bonjour")).toEqual("Bonjour Luke Skywalker, may the force be with you!");
        expect(result.current.jediGreeting()).toEqual("Hello Luke Skywalker, may the force be with you!");
    });

    it("should return stable callback references when state and factoryArgs keeps the same reference", () => {
        const factoryArgs = {force: "Force"};

        const {result, rerender} = renderHook(
            ({state, factoryArgs}) =>
                useThunks(declarations, factoryArgs, state),
            {
                initialProps: {
                    state: initialState,
                    factoryArgs: factoryArgs
                }
            }
        );
        const callbacks1 = result.current;
        rerender({state: initialState, factoryArgs: factoryArgs});
        const callbacks2 = result.current;
        expect(callbacks1.jediGreeting).toBe(callbacks2.jediGreeting);
        expect(callbacks1.sithGreeting).toBe(callbacks2.sithGreeting);
    });

    it("should return stable callback references when factoryArgs keeps the same reference, even if the state changes as long as what selected doesn't", () => {
        const factoryArgs = {force: "Force"};

        const {result, rerender} = renderHook(
            ({state, factoryArgs}) =>
                useThunks(declarations, factoryArgs, state),
            {
                initialProps: {
                    state: initialState,
                    factoryArgs: factoryArgs
                }
            }
        );

        const callbacks1 = result.current;
        rerender({state: {...initialState, jedi: [...initialState.jedi, {name: "Obi Wan Kenobi", lightsaber: "blue"}]}, factoryArgs: factoryArgs});
        const callbacks2 = result.current;

        expect(callbacks1.jediGreeting).toBe(callbacks2.jediGreeting);
        expect(callbacks1.sithGreeting).toBe(callbacks2.sithGreeting);
    });

    it("should update the callback reference when the partial state changes", () => {
        const factoryArgs = {force: "Force"};

        const {result, rerender} = renderHook(
            ({state, factoryArgs}) =>
                useThunks(declarations, factoryArgs, state),
            {
                initialProps: {
                    state: initialState,
                    factoryArgs: factoryArgs
                }
            }
        );
        const callbacks1 = result.current;
        // Change the first jedi (partial state for jediGreeting)
        const newState = {
            ...initialState,
            jedi: [{name: "Yoda", lightsaber: "green"}]
        };
        rerender({state: newState, factoryArgs: factoryArgs});
        const callbacks2 = result.current;

        expect(callbacks2.jediGreeting("Bonjour")).toEqual("Bonjour Yoda, may the force be with you!");

        expect(callbacks1.jediGreeting).not.toBe(callbacks2.jediGreeting);
        // sithGreeting should remain the same
        expect(callbacks1.sithGreeting).toBe(callbacks2.sithGreeting);

    });
});