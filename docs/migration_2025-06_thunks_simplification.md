# Migration: Thunk Simplification and Aggregation (2025-06)

This document outlines a refactoring task for LLM agents to simplify thunk declarations and improve their organization within the codebase. The goal is to leverage `createSimpleThunkDeclaration` where applicable and aggregate related thunks into a single exported object within their respective feature directories.

## Task Description

For a given feature directory (e.g., `web/src/core/catalog/album-create`), perform the following steps:

1.  **Simplify Thunk Declarations**:
    *   Review each `ThunkDeclaration` within the target directory (e.g., `web/src/core/catalog/album-create/thunk-createAlbum.ts`).
    *   If a thunk's `factory` function solely dispatches a single action (i.e., its logic is `({dispatch}) => (...args) => { dispatch(actionCreator(...args)); }`), simplify its declaration by replacing the explicit `ThunkDeclaration` object with a call to `createSimpleThunkDeclaration`.
    *   **Example Transformation**:

        **Before:**
        ```typescript
        // web/src/core/catalog/album-create/thunk-createAlbum.ts
        import {createAlbumAction} from "./action-createAlbum";
        import {ThunkDeclaration} from "src/libs/dthunks";

        export const createAlbumDeclaration: ThunkDeclaration<
            CatalogViewerState,
            {},
            (name: string) => void,
            CatalogFactoryArgs
        > = {
            selector: () => ({}),
            factory: ({dispatch}) => {
                return (name: string) => {
                    dispatch(createAlbumAction(name));
                };
            },
        };
        ```

        **After:**
        ```typescript
        // web/src/core/catalog/album-create/thunk-createAlbum.ts
        import {createAlbumAction} from "./action-createAlbum";
        import {createSimpleThunkDeclaration} from "src/libs/dthunks";

        export const createAlbumDeclaration = createSimpleThunkDeclaration(createAlbumAction);
        ```
    *   **Important**: If the `factory` contains additional logic (e.g., API calls, conditional dispatching, complex state manipulation), *do not* simplify it with `createSimpleThunkDeclaration`. Leave it as an explicit `ThunkDeclaration` object.

2.  **Aggregate Thunks into a Single Export**:
    *   In the `index.ts` file of the target directory (e.g., `web/src/core/catalog/album-create/index.ts`), create a new `const` export named `[featureName]Thunks` (e.g., `albumCreateThunks`).
    * Remove thunks exports
    * Import the thunks to be used
    *   This new object should contain all the `ThunkDeclaration` exports from the directory.
    *   **Example**:

        **Before (`web/src/core/catalog/album-create/index.ts`):**
        ```typescript
        export * from "./action-createAlbum";
        export * from "./thunk-createAlbum";
        // ... other exports
        ```

        **After (`web/src/core/catalog/album-create/index.ts`):**
        ```typescript
        // Removes the thunks exports
        // Keep the selectors and action exports
        import {createAlbumDeclaration} from "./thunk-createAlbum" // add the required imports 

        /**
         * Thunks related to album creation.
         *
         * Expected handler types:
         * - `createAlbum`: `(name: string) => void`
         */
        export const albumCreateThunks = {
            createAlbum: createAlbumDeclaration,
        };
        ```

3.  **Update Main `thunks.ts` Aggregation**:
    *   In `web/src/core/catalog/thunks.ts`, import the newly created `[featureName]Thunks` object.
    *   Remove the individual imports for the thunk declarations that are now part of the aggregated object.
    *   Spread the `[featureName]Thunks` object into the main `catalogThunks` object.
    *   **Example**:

        **Before (`web/src/core/catalog/thunks.ts`):**
        ```typescript
        import {createAlbumDeclaration} from "./album-create";
        // ... other imports

        export const catalogThunks = {
            createAlbum: createAlbumDeclaration,
            // ... other thunks
        };
        ```

        **After (`web/src/core/catalog/thunks.ts`):**
        ```typescript
        import {albumCreateThunks} from "./album-create"; // New import
        // ... other imports, remove individual thunks import

        export const catalogThunks = {
            ...albumCreateThunks, // Spread the new aggregated object
            // ... other thunks
        };
        ```

4.  **Add JSDoc Documentation to Aggregated Thunks**:
    *   Add a JSDoc comment to the newly created `[featureName]Thunks` object in the `index.ts` file.
    *   This comment should include a brief description and a list of `Expected handler types` for each thunk within the object. This helps LLM agents understand the expected function signature when these thunks are used as handlers.
    *   **Example**:
        ```typescript
        /**
         * Thunks related to album creation.
         *
         * Expected handler types:
         * - `createAlbum`: `(name: string) => void`
         */
        export const albumCreateThunks = {
            createAlbum: createAlbumDeclaration,
        };
        ```
