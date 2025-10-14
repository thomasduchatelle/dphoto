# Coding Principles

You are a strong developer prioritizing simple and well tested code. You **strictly follow the coding principles** defined below.

1. IMPORTANT: **do no add comments in your code**. ONLY use the chat to communicate, NEVER the file listing.
2. Prioritise exact semantic over correct syntax: syntax error can be found by the compiler and fixed easily, semantic error could void a test without being noticed.
3. Avoid updating existing tests unless explicitly requested: create new ones.
    * if a test would fail otherwise, update it and bring it to the user attention in the chat: "WARNING - TEST
      UPDATED: <name of the test and reason to be updated>"
    * if a new test makes another one redundant, bring it to the user attention in the chat: "INFO - TEST REDUNDANT: <name of the existing test> -> name of the
      new tests"
4. Always use explicit types or interfaces, never use `any`.
5. Inline the variables where possible.

## Testing Strategy

Components holding a logic are tested to demonstrate the Acceptance Criteria are fulfilled. All the tests must follow these principles:

1. **TDD principle**: implementations should **never have a behavior that hasn't been expected or forced by a test case**. Without an appropriate test, code
   must remain extremely simple, even if it means it is wrong.

2. **Robust Tests**: tests must be robust to refactoring
    * actions and selectors are **tested together**: initiate the state -> create action -> execute the reducer -> run the selector -> assert the returned
      selection
    * use pre-defined constants from `web/src/core/catalog/tests/test-helper-state.ts`
    * use fake implementations: method signature of the Ports can change, the fakes are updated, but it does not affect the tests

3. **Unit test first**: the code is structured so the complex logic of acceptance criteria can be validated by unit tests. Use the appropriate framework:
    * actions and selectors: jest unit tests
    * thunks: jest unit tests
    * UI components: Ladle (Storybook like) with the component set in each of the relevant situations (examples: "default", "loading", "with a technical error", "success")
    * react hooks: jest + testing-library

## Domain-Specific References

When working on different domains, replace the example types and adapters with the appropriate ones for your domain:

**Catalog Domain:**
- State: `CatalogViewerState` 
- Factory Args: `CatalogFactoryArgs`
- Adapter: `new CatalogFactory(app as DPhotoApplication).restAdapter()`

**Other Domains:**
- State: `{DomainName}State` (replace `{DomainName}` with your actual domain name)
- Factory Args: `{DomainName}FactoryArgs`
- Adapter: Follow the specific adapter pattern for your domain

**Note**: The examples in this document use `StarWars` as the domain name for educational purposes. Replace `StarWars` with your actual domain name (e.g.,
`Catalog`, `Security`, etc.) when implementing.

## Tree structure

The file structure is as follows:

* `components/`
    * `catalog-react/` - contains the React Components used to integrate the domain to the other components
    * `EditDateDialog` - ... as an example, the list of the components
* `core/catalog/` - "catalog" is the name of the domain
    * `<feature name in dash-case>/` - each feature has a folder containing related actions, thunks, and selectors
    * `language/` - ubiquitous language and definition of the State shared for the domain
    * `common/` - functionalities reused across most features
    * `actions.ts` - where the action interface and the partial reducer are registered
    * `thunks.ts` - where the thunks are registered
    * `adapters/api/` - where the REST API adapter are implemented

## Design Principles

In the UI, the data flow through:

1. **UI Components (Views)**: purely presentational. They receive all necessary data via props (derived from the Redux state via selectors) and trigger
   thunks in response to user interactions. They do not contain any business logic or direct state manipulation.

2. **Thunks**: encapsulate all business logic. They are triggered by UI events, interact with adapters (ports) for external
   operations (like API calls), and dispatch actions to update the application state. They will be stateless and framework-agnostic.
    * **Port** - an interface declared next to the Thunk to abstract specific technologies (REST API, Local store, ...)
    * **Adapter** - a class implementation of the _Port_ for a specific technology (example: Axios)

3. **Actions**: sole mechanism for state mutation. Each action represent a specific event in the system

4. **State**: single source of truth for the UI, it uses ubiquitous language

5. **Selectors**: responsible for transforming raw state data into the specific data structures required by the UI components

When designing or refactoring, apply **strictly** the following rules:

1. Only add on the UI components the properties they **require** for the **current story**
    * All data properties must come from a selector
    * All callback function must come from a thunk

2. Keep the payload of action and thunk **minimum**:
    * _actions_ payload is only what's new on the state, or the Identifiers required to mutate it
    * _thunks_ arguments is only what it is required to perform the business logic (REST requests, making decisions, ...) ; only extract from the state what's
      needed ; only take as parameters what cannot be derived from the state

## Naming Conventions and Coding Practices

### Selectors

Selector function naming convention is `<name in camelCase>Selector`, and they return a selection with the naming convention `<name in PascalCase>Selection`.

Example:

```typescript
function jediAcademySelector(state: StarWarsState): JediAcademySelection {
}
```

### Actions

The _Action_ naming convention is the **Past Tense Event Naming** and is represented in **camelCase** (example: `albumCreated`)

An _Action_ is defined using the `createAction` function in a single file named `.../<feature in dash-case>/action-<action name in camelCase>.ts`.

Here a good example on how to test an action:

```typescript
// web/src/core/starwars/language/state.ts
interface ForceWielder {
    name: string
    lightsaber: string
}

interface StarWarsState {
    jedi: ForceWielder[]
    sith: ForceWielder[]
    balanceInTheForce: "troubled" | "balanced"
}

// web/src/core/starwars/tests/test-helper-state.ts
const stateAtEpisode6: StarWarsState = {
    jedi: [{name: "Luke Skywalker", lightsaber: "Green"}],
    sith: [{name: "Palpatine", lightsaber: "Red"}, {name: "Darth Vader", lightsaber: "Red"}],
    balanceInTheForce: "troubled",
};

const jediAcademyAtEpisode6: JediAcademySelection = {
    master: {name: "Luke Skywalker", lightsaber: "Green"},
    balanceInTheForce: "troubled",
}

// web/src/core/starwars/battles/selector-jediAcademySelector.ts
import {StarWarsState, ForceWielder} from "../language/state";

export interface JediAcademySelection {
    master: ForceWielder
    balanceInTheForce: "troubled" | "balanced"
}

export function jediAcademySelector(state: StarWarsState): JediAcademySelection {
    // .... implementation skipped
};

// web/src/core/starwars/battles/action-forceWielderDefeated.test.ts
import {forceWielderDefeated} from "./action-forceWielderDefeated";
import {jediAcademySelector, JediAcademySelection} from "./selector-jediAcademySelector";
import {StarWarsState} from "../language/state";
import {stateAtEpisode6, jediAcademyAtEpisode6} from "../tests/test-helper-state";

describe('action:forceWielderDefeated', () => { // naming convention of "describe" is `action:<action name>`

    const sithLordsAtEpisode6 = ["Palpatine", "Darth Vader"]
    const forceIsBalanced = "balanced"

    it('should balance the force if all Siths have been defeated', () => {
        // 1. Initiate a state
        const state = {
            ...stateAtEpisode6,
            // specific properties can be overridden for this test purpose
        };

        // 2. Execute the reducer
        const action = forceWielderDefeated(sithLordsAtEpisode6);
        const got = action.reducer(state, action);

        // 3. Execute the selector
        // 4. Assert the whole result of the selector
        expect(jediAcademySelector(got)).toEqual({
            ...jediAcademyAtEpisode6,   // assert that properties that shouldn't have changed haven't changed (using a constant)
            balanceInTheForce: forceIsBalanced, // explicitly validate the properties expected to have changed ; use const with explicit names if it helps readability
        });
    });
});
```

#### Case 1: No payload

```typescript
// web/src/core/starwars/battles/action-battleStarted.ts
import {createAction} from "src/libs/daction";
import {StarWarsState} from "src/core/starwars/language/state";

export const battleStarted = createAction<StarWarsState>(
    "BattleStarted",
    (current: StarWarsState) => {
        return {
            ...current,
            balanceInTheForce: "troubled",
        };
    }
);

export type BattleStarted = ReturnType<typeof battleStarted>;
```

#### Case 2: Single payload

```typescript
// web/src/core/starwars/battles/action-forceWielderDefeated.ts
import {createAction} from "src/libs/daction";
import {StarWarsState} from "src/core/starwars/language/state";

export const forceWielderDefeated = createAction<StarWarsState, string[]>(
    "ForceWielderDefeated",
    (current: StarWarsState, defeatedWielderNames: string[]) => {
        const updatedJedi = current.jedi.filter(jedi => !defeatedWielderNames.includes(jedi.name));
        const updatedSith = current.sith.filter(sith => !defeatedWielderNames.includes(sith.name));

        return {
            ...current,
            jedi: updatedJedi,
            sith: updatedSith,
        };
    }
);

export type ForceWielderDefeated = ReturnType<typeof forceWielderDefeated>;
```

#### Case 3: Multiple properties

```typescript
// web/src/core/starwars/battles/action-jediTurned.ts
import {createAction} from "src/libs/daction";
import {StarWarsState, ForceWielder} from "src/core/starwars/language/state";

interface JediTurnedPayload {
    jediName: string;
    sithLordName: string
}

export const jediTurned = createAction<StarWarsState, JediTurnedPayload>(
    "JediTurned",
    ({jedi, sith, ...current}: StarWarsState, {jediName, sithLordName}: JediTurnedPayload) => {
        return {
            ...current,
            jedi: jedi.filter(j => j.name !== jediName),
            sith: [...sith, {name: `Darth ${sithLordName}`, lightsaber: "Red"}]
        }
    }
);

export type JediTurned = ReturnType<typeof jediTurned>;
```

### Thunks

The thunk naming convention is a verb with a complement, example: `createAlbum`, or `sendEmail`

#### Case 1: complete example with: port, preselection, and argument

```typescript
// web/src/core/starwars/battles/thunk-executeOrder66.ts
import {ThunkDeclaration} from "src/libs/dthunks";
import {StarWarsState} from "../language/state";
import {forceWielderDefeated, ForceWielderDefeated} from "./action-forceWielderDefeated";
import {StarWarsFactoryArgs} from "../common/starwars-factory-args";

// Port naming convention is "<thunk name>Port"
export interface ExecuteOrder66Port {
    eliminateJedi(orderId: string, jediNames: string[]): Promise<string[]>; // Now accepts orderId and jediNames
}

// The business logic is implemented in the thunk function named "<thunk name>Thunk"
export async function executeOrder66Thunk(
    dispatch: (action: ForceWielderDefeated) => void,   // first argument is always "dispatch" accepting the union of dispatched actions.
    order66Port: ExecuteOrder66Port,   // (optional) port used by the thunk
    jediToEliminate: string[],  // (optional) pre-selection and others bound by the factory
    orderId: string             // (optional) last argument(s) are the values passed from the UI
): Promise<void> {
    const defeatedWielders = await order66Port.eliminateJedi(orderId, jediToEliminate);
    dispatch(forceWielderDefeated(defeatedWielders));
}

export const executeOrder66Declaration: ThunkDeclaration<
    StarWarsState, // this is the global state of the domain
    { jediToEliminate: string[] }, // type of the pre-selection from the global state
    (orderId: string) => Promise<void>, // function visible by the UI, and returned by the factory
    StarWarsFactoryArgs // constant of the domain {app: DPhotoApplication ; dispatch: (action: Action<StarWarsState, any>) => void }
> = {
    // Pre-selector: extracts relevant values from the state
    selector: ({jedi}: StarWarsState) => ({
        jediToEliminate: jedi.map(jedi => jedi.name),
    }),

    // Factory: wires up dependencies and returns the handler used by UI components
    factory: ({dispatch, app, partialState: {jediToEliminate}}) => {
        // get or create the adapter, exmaple for catalog domain:
        //  const updateAlbumDatesPort: UpdateAlbumDatesPort = new CatalogAPIAdapter(app.axiosInstance, app);
        const order66Port: ExecuteOrder66Port = app.starwarsAdapter();
        return executeOrder66Thunk.bind(null, dispatch, order66Port, jediToEliminate);
    },
};
```

#### Case 2: thunk solely dispatches a single action

```typescript
// web/src/core/starwars/battles/thunk-acceptPadawan.ts
import {padawanAccepted} from "./action-padawanAccepted";
import {createSimpleThunkDeclaration} from "src/libs/dthunks";
import {StarWarsState} from "../language/state";

export const acceptPadawanDeclaration = createSimpleThunkDeclaration(padawanAccepted);
```

#### Make the thunk available in UI

```typescript
// web/src/core/starwars/battles/index.ts
import {acceptPadawanDeclaration} from "./thunk-acceptPadawan"
import {executeOrder66Declaration} from "./thunk-executeOrder66"

// the following JSDoc must be maintained: the file is used by LLM agents to know what handlers are available. This is the only exception to the rule "do not add comment".

/**
 * Starwars' battles feature exposes the handlers:
 *
 * - `acceptPadawan`: `(padawan: ForceWielder) => void`
 * - `executeOrder66`: `(orderId: string) => Promise<void>`
 */
export const battlesThunks = {
    acceptPadawan: acceptPadawanDeclaration,
    executeOrder66: executeOrder66Declaration,
};


export const thunks = {
    ...battlesThunks,
    // other aggregated thunks from other features
}
```

#### Testing with Fake

Do not use MOCKs nor STUBS. **You must use FAKE implementations as follow**:

```typescript
// web/src/core/starwars/battles/thunk-executeOrder66.test.ts
import {executeOrder66Thunk, ExecuteOrder66Port} from "./thunk-executeOrder66";
import {forceWielderDefeated} from "./action-forceWielderDefeated";
import {StarWarsState} from "../language/state";
import {Action} from "src/libs/daction";

// Fake implementation reproduce the expected behaviour of the actual implementation
class BattleFakeAdapter implements ExecuteOrder66Port {
    constructor(
        private jediDyingByOrder: Map<string, string[]>
    ) {
    }

    public eliminateJedi(orderId: string, jediNames: string[]): Promise<string[]> {
        return Promise.resolve(jediNames.filter(jedi => {
            const orders = this.jediDyingByOrder.get(jedi) ?? []
            return orders.indexOf(orderId) >= 0
        }));
    }
}

describe("thunk:executeOrder66Thunk", () => {
    it("should dispatch forceWielderDefeated with the list of defeated Jedi from the port", async () => {
        const order66Id = "Order66";
        const fakePort = new BattleFakeAdapter(new Map([
            ["Obi-Wan Kenobi", []],
            ["Aayla Secura", [order66Id]],
        ]));

        const jediToEliminate = ["Obi-Wan Kenobi", "Yoda", "Aayla Secura"];
        const dispatchedActions: Action<StarWarsState, any>[] = [];

        await executeOrder66Thunk(dispatchedActions.push.bind(dispatchedActions), fakePort, jediToEliminate, order66Id);

        expect(dispatchedActions).toEqual([
            forceWielderDefeated(["Aayla Secura"]),
        ]);
    });
});
```

### UI Component testing: Ladle Stories

UI components must be tested visually, using Ladle, validating each of their relevant situations (default, saving, displaying an error, loading, ...).

#### Example for a simple component (5 properties or fewer):

```typescript jsx
// album-card.stories.tsx
import {action} from "@ladle/react";

export const Default = <AlbumCard name='Jan 2025' size='42' onClick={action('onClick')}/>
export const Disabled = <AlbumCard name='Jan 2025' size='42' disabled/>
```

#### Example for complex components (more than 5 properties, or openable components like dialogs)

This rules must be followed to make the tests easy to read and update:

1. a test component must be created to wrap the component under tests. It is named with `Wraper` suffix (example: `DeleteAlbumDialogWrapper` to wrap `DeleteAlbumDialog`)
2. a React state must be created when a callback match with a property (example: `onNameChanged` callback and `name` property)
3. `action('callbackName')` must be used for the other callbacks (example: `<DeleteAlbumDialog onSubmit={action('onSubmit')}`)
4. if a property `open` is present (or one with a similar meaning):
   * its value must be managed by a state, defaulted to `true`. Example: `const [open, setOpen] = useState(true)`
   * a "Reopen Dialog" button must be present to set the value to `true`
5. other properties can be overridden by setting the args of the stary (example: `Disabled.args = { name: 'Empire Strike Back' }`)

A complete example of the stories of the component `export const DeleteAlbumDialog = (args: DeleteAlbumDialogProps) => { ... }`:

```typescript jsx
// delete-album-dialog.stories.tsx
import {action, Story} from "@ladle/react";

type Props = Partial<DeleteAlbumDialogProps>

const DeleteAlbumDialogWrapper : Story<Props> = (props: Props) => {
    const [open, setOpen] = React.useState(true);
    const [albumName, setAlbumName] = React.useState(props.albumName);

    return (
        <>
            <Button variant='contained' onClick={() => setOpen(true)}>
                Reopen Dialog
            </Button>
            <DeleteAlbumDialog
                albumId='sensible-default-id'
                {...props}
                open={open}
                albumName={albumName}
                onClose={() => setOpen(false)}
                onDelete={action("onDelete")}
            />
        </>
    );
};

export const Default = (args) => <DeleteAlbumDialogWrapper {...args} />
Default.args = {albumId: 'episode-5', albumName: 'Empire Strike Back'}
```
