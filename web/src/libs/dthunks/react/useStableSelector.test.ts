import {renderHook} from "@testing-library/react";

import {useStableSelector} from "./useStableSelector";

describe("useStableSelector", () => {
    // Complex Star Wars state
    const state = {
        jedi: [
            {name: "Luke Skywalker", lightsaber: "green", master: "Obi-Wan Kenobi"},
            {name: "Obi-Wan Kenobi", lightsaber: "blue", master: "Qui-Gon Jinn"}
        ],
        sith: [
            {name: "Darth Vader", lightsaber: "red", master: "Emperor Palpatine"},
            {name: "Darth Maul", lightsaber: "red", master: "Darth Sidious"}
        ],
        droids: [
            {name: "R2-D2", type: "astromech"},
            {name: "C-3PO", type: "protocol"}
        ],
        planets: {
            tatooine: {climate: "arid", population: "200000"},
            dagobah: {climate: "murky", population: "unknown"}
        }
    };

    const firstJediSelector = (input: typeof state) => input.jedi[0];


    function createStableSelector(selector: any): typeof selector {
        const {result} = renderHook(() =>
            useStableSelector(selector)
        );
        return result.current;
    }

    it("should return the same selected result than the original selector", () => {
        // Use the hook
        const stableSelector = createStableSelector(firstJediSelector);

        expect(stableSelector(state)).toEqual(firstJediSelector(state));
    });

    it("should return the same reference of a function when the component is refreshed with the same selector", () => {
        // Simulate a component refresh by re-rendering the hook with the same selector
        const {result, rerender} = renderHook(() => useStableSelector(firstJediSelector));
        const stableSelector1 = result.current;
        rerender();
        const stableSelector2 = result.current;
        expect(stableSelector1).toBe(stableSelector2);
    });

    it("should return twice the same reference when called with the state, and a modified state that has the same first jedi (same reference as well)", () => {
        const stableSelector = createStableSelector(firstJediSelector);

        const result1 = stableSelector(state);

        // Create a new state object, but reuse the same jedi array (same reference for first jedi)
        const modifiedState = {
            ...state,
            sith: [...state.sith], // new array, but irrelevant
            droids: [...state.droids],
            planets: {...state.planets}
        };

        const result2 = stableSelector(modifiedState);

        expect(result1).toBe(result2);
    });

    it("should return a the same reference if the selector return the same values (shallow equal)", () => {
        const stableSelector = createStableSelector(firstJediSelector);

        const result1 = stableSelector(state);

        // Create a new state with a new first jedi object (different reference)
        const newFirstJedi = {name: "Luke Skywalker", lightsaber: "green", master: "Obi-Wan Kenobi"};
        const updatedState = {
            ...state,
            jedi: [{...newFirstJedi}, ...state.jedi.slice(1)]
        };

        const result2 = stableSelector(updatedState);

        expect(result1).toBe(result2);
        expect(result2).toEqual(newFirstJedi);
    });

    it("should return a different reference if the selector returns something different", () => {
        const stableSelector = createStableSelector(firstJediSelector);

        const result1 = stableSelector(state);

        // Create a new state with a new first jedi object (different reference)
        const newFirstJedi = {name: "Yoda", lightsaber: "green", master: "The Force"};
        const updatedState = {
            ...state,
            jedi: [{...newFirstJedi}, ...state.jedi.slice(1)]
        };

        const result2 = stableSelector(updatedState);

        expect(result1).toEqual({name: "Luke Skywalker", lightsaber: "green", master: "Obi-Wan Kenobi"});
        expect(result2).toEqual({name: "Yoda", lightsaber: "green", master: "The Force"});

        expect(result1).not.toBe(result2);
        expect(result2).toEqual(newFirstJedi);
    });
});
