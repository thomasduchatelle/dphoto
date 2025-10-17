---
applyTo: "web/**"
---

# Coding Principles

You are a strong developer prioritizing simple and well tested code. You **strictly follow the coding principles** defined below.

1. IMPORTANT: **do no add comments in your code**. ONLY use the chat to communicate, NEVER the file listing.
2. Prioritise exact semantic over correct syntax: syntax error can be found by the compiler and fixed easily, semantic error could void a test without being
   noticed.
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
    * UI components: Ladle (Storybook like) with the component set in each of the relevant situations (examples: "default", "loading", "with a technical
      error", "success")
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

**File Convention**: Each selector is in its own file `selector-<name>.ts` to keep business logic modular and discoverable.

**Export Format**: Always export selectors as named exports (not default). This enables better IDE support, refactoring, and tree-shaking.

```typescript
// selector-availableJedi.ts
export const selectAvailableJedi = (state: StarWarsState): Jedi[] => {
    return state.jedi.filter(j => j.status === 'alive');
};
```

**Naming Pattern**: Use the prefix `select` followed by the data being selected. The name should clearly indicate what data will be returned.

```typescript
// Good Examples
export const selectCurrentUser = (state: AppState): User => { ...
};
export const selectUnreadMessages = (state: AppState): Message[] => { ...
};
export const selectIsLoadingProfile = (state: AppState): boolean => { ...
};

// Poor Examples (avoid these)
export const currentUser = (state: AppState): User => { ...
};     // Missing 'select' prefix
export const getUser = (state: AppState): User => { ...
};         // Inconsistent verb
```

**Testing Pattern**: Always test selectors with their corresponding actions to ensure state transformations work as expected:

```typescript
// selector-availableJedi.test.ts
describe('selectAvailableJedi', () => {
    it('should return only alive Jedi after Order 66', () => {
        const initialState: StarWarsState = {jedi: ALL_JEDI};

        const action = forceWielderDefeated(['Aayla Secura', 'Ki-Adi-Mundi']);
        const newState = reducer(initialState, action);

        const result = selectAvailableJedi(newState);

        expect(result).toEqual([
            {name: 'Obi-Wan Kenobi', status: 'alive'},
            {name: 'Yoda', status: 'alive'}
        ]);
    });
});
```

### Actions

**File Convention**: One action per file named `action-<eventName>.ts` where the event name is in camelCase (not PascalCase).

**Naming Pattern**: Action creator functions use camelCase and describe the **event that occurred** (not the desired outcome):

```typescript
// Good Examples (event-driven)
export const userLoggedIn = (userId: string) => ({...});
export const albumDeleted = (albumId: string) => ({...});
export const mediaUploadFailed = (error: Error) => ({...});

// Poor Examples (avoid these)
export const setUser = (userId: string) => ({...});           // Sounds like a command
export const deleteAlbum = (albumId: string) => ({...});      // Sounds like a command
export const UserLoggedIn = (userId: string) => ({...});      // Wrong case (should be camelCase)
```

**Factory Function Signature**:

```typescript
// action-userLoggedIn.ts
import {Action} from 'src/libs/daction';
import {AppState} from '../language/state';

export const userLoggedIn = (userId: string): Action<AppState, { userId: string }> => ({
    type: 'userLoggedIn',
    payload: {userId},
    reduce: (state) => ({
        ...state,
        currentUserId: userId,
        isAuthenticated: true
    })
});
```

### Thunks

**File Convention**: Each thunk is in its own file `thunk-<operation>.ts`, named after the business operation it performs.

**Naming Pattern**: Thunk names should use camelCase and describe the business operation, suffixed with `Thunk`:

```typescript
// Good Examples
export const loadUserProfileThunk = async (...) => { ...
};
export const deleteAlbumThunk = async (...) => { ...
};
export const uploadMediaThunk = async (...) => { ...
};

// Poor Examples
export const loadUserProfile = async (...) => { ...
};   // Missing 'Thunk' suffix
export const LoadUserProfileThunk = async (...) => { ...
}; // Wrong case
```

**Thunk Signature Pattern**:

```typescript
// thunk-deleteAlbum.ts
export interface DeleteAlbumPort {
    deleteAlbum(albumId: string): Promise<void>;
}

export const deleteAlbumThunk = async (
    dispatch: Dispatch<AppState>,
    port: DeleteAlbumPort,
    albumId: string
): Promise<void> => {
    dispatch(albumDeletionStarted(albumId));

    try {
        await port.deleteAlbum(albumId);
        dispatch(albumDeleted(albumId));
    } catch (error) {
        dispatch(albumDeletionFailed(albumId, error));
    }
};
```

**Testing Thunks with Fakes**: Use fake implementations instead of mocks. Fakes should mimic real behavior and be reusable:

```typescript
// thunk-executeOrder66.test.ts
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

1. a test component must be created to wrap the component under tests. It is named with `Wraper` suffix (example: `DeleteAlbumDialogWrapper` to wrap
   `DeleteAlbumDialog`)
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

const DeleteAlbumDialogWrapper: Story<Props> = (props: Props) => {
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
