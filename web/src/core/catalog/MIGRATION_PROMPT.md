# Action Migration Prompt

You are tasked with migrating Redux-style actions from the old verbose pattern to the new simplified `createAction` pattern. This migration reduces boilerplate and eliminates the need for separate interfaces, factory functions, and reducer registrations.

## Overview

**OLD PATTERN** (verbose):
```typescript
// Interface
export interface SomeAction {
    type: "SomeAction";
    payload?: SomePayload;
}

// Factory function
export function someAction(payload?: SomePayload): SomeAction {
    return { type: "SomeAction", payload };
}

// Reducer function
export function reduceSomeAction(state: State, action: SomeAction): State {
    // reducer logic
    return newState;
}

// Registration function
export function someActionReducerRegistration(handlers: any) {
    handlers["SomeAction"] = reduceSomeAction;
}
```

**NEW PATTERN** (simplified):
```typescript
import { createAction } from "../common/action-factory";

export const someAction = createAction<State, PayloadType>(
    "SomeAction",
    (state: State, payload: PayloadType) => {
        // reducer logic
        return newState;
    }
);

export type SomeAction = ReturnType<typeof someAction>;
```

## Migration Steps

### 1. Identify Action Patterns

Look for files with this structure:
- `action-*.ts` files containing interfaces, factory functions, and reducers
- Corresponding `action-*.test.ts` files
- Registration functions ending with `ReducerRegistration`

### 2. Transform Actions

#### For actions WITHOUT payload:
```typescript
// OLD
export interface ActionName {
    type: "ActionName";
}

export function actionName(): ActionName {
    return { type: "ActionName" };
}

export function reduceActionName(state: State, action: ActionName): State {
    // logic here
    return newState;
}

// NEW
export const actionName = createAction<State>(
    "ActionName",
    (state: State) => {
        // logic here
        return newState;
    }
);
```

#### For actions WITH single payload:
```typescript
// OLD
export interface ActionName {
    type: "ActionName";
    payload: PayloadType;
}

export function actionName(payload: PayloadType): ActionName {
    return { type: "ActionName", payload };
}

export function reduceActionName(state: State, {payload}: ActionName): State {
    // logic using payload
    return newState;
}

// NEW
export const actionName = createAction<State, PayloadType>(
    "ActionName",
    (state: State, payload: PayloadType) => {
        // logic using payload
        return newState;
    }
);
```

#### For actions WITH multiple properties:
```typescript
// OLD
export interface ActionName {
    type: "ActionName";
    prop1: Type1;
    prop2: Type2;
}

export function actionName(props: Omit<ActionName, "type">): ActionName {
    return { ...props, type: "ActionName" };
}

export function reduceActionName(state: State, {prop1, prop2}: ActionName): State {
    // logic using prop1, prop2
    return newState;
}

// NEW
interface ActionNamePayload {
    prop1: Type1;
    prop2: Type2;
}

export const actionName = createAction<State, ActionNamePayload>(
    "ActionName",
    (state: State, {prop1, prop2}: ActionNamePayload) => {
        // logic using prop1, prop2
        return newState;
    }
);
```

### 3. Update Tests

Transform test files to use the new action pattern:

```typescript
// OLD
import { actionName, reduceActionName } from "./action-actionName";

it("test description", () => {
    const action = actionName(payload);
    const result = reduceActionName(state, action);
    // assertions
});

// NEW
import { actionName } from "./action-actionName";

it("test description", () => {
    const action = actionName(payload);
    const result = action.reducer(state, action);
    // assertions
});

// Add comparison tests
it("supports action comparison for testing", () => {
    const action1 = actionName(payload);
    const action2 = actionName(payload);
    
    expect(action1).toEqual(action2);
    expect([action1]).toContainEqual(action2);
});
```

### 4. Update Main Actions File

Remove migrated actions from the registration system:

```typescript
// Remove imports of old reducer registration functions
// Remove from reducerRegistrations array
// Remove from CatalogViewerAction union type
// Remove from export type declarations
```

### 5. Always Include

- Import `createAction` from `"../common/action-factory"`
- Export the type: `export type ActionName = ReturnType<typeof actionName>;`
- Update all test files to use `action.reducer(state, action)` instead of separate reducer functions
- Add action comparison tests to verify testing compatibility

## Example Migration

**Before:**
```typescript
// action-userLoggedIn.ts
export interface UserLoggedIn {
    type: "UserLoggedIn";
    userId: string;
    username: string;
}

export function userLoggedIn(userId: string, username: string): UserLoggedIn {
    return { type: "UserLoggedIn", userId, username };
}

export function reduceUserLoggedIn(state: AppState, {userId, username}: UserLoggedIn): AppState {
    return {
        ...state,
        currentUser: { userId, username },
        isAuthenticated: true,
    };
}

export function userLoggedInReducerRegistration(handlers: any) {
    handlers["UserLoggedIn"] = reduceUserLoggedIn;
}
```

**After:**
```typescript
// action-userLoggedIn.ts
import { createAction } from "../common/action-factory";

interface UserLoggedInPayload {
    userId: string;
    username: string;
}

export const userLoggedIn = createAction<AppState, UserLoggedInPayload>(
    "UserLoggedIn",
    (state: AppState, {userId, username}: UserLoggedInPayload) => {
        return {
            ...state,
            currentUser: { userId, username },
            isAuthenticated: true,
        };
    }
);

export type UserLoggedIn = ReturnType<typeof userLoggedIn>;
```

## Key Benefits

1. **Reduced boilerplate**: Single declaration instead of 4+ separate pieces
2. **Type safety**: Full TypeScript inference and checking
3. **Auto-registration**: No manual reducer registration needed
4. **Testing compatibility**: Actions remain comparable with `toEqual()`
5. **Backward compatibility**: Works alongside existing legacy actions

## Important Notes

- The `createAction` function supports 0-3 parameters automatically
- For tuple payloads, use: `createAction<State, [Type1, Type2]>`
- Always test action comparison to ensure testing compatibility
- The generic reducer automatically handles actions with built-in reducers
- Legacy actions continue to work during the migration period

Migrate one folder at a time, ensuring all tests pass before moving to the next folder.
