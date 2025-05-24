import {CatalogViewerState} from "./catalog-state";
import {reduceAlbumsAndMediasLoaded} from "./action-albumsAndMediasLoadedAction";
import {reduceAlbumsLoaded} from "./action-albumsLoadedAction";
import {reduceMediaFailedToLoad} from "./catalog-action-MediaFailedToLoadAction";
import {CatalogSupportedActions} from "./catalog-index";

/**
 * Utility to create a reducer from a map of action handlers.
 * 
 * Usage:
 *   const reducer = createReducer({
 *     AlbumsLoadedAction: (state, action) => { ... },
 *     MediasLoadedAction: (state, action) => { ... },
 *   });
 * 
 * The returned reducer takes (state, action) and dispatches to the correct handler.
 */
export function createReducer<TState, TActions extends { type: string }>(
    handlers: {
        [K in TActions["type"]]: (state: TState, action: Extract<TActions, { type: K }>) => TState
    }
): (state: TState, action: TActions) => TState {
    return (state: TState, action: TActions): TState => {
        const handler = handlers[action.type as keyof typeof handlers];
        if (handler) {
            // TypeScript will ensure correct type for handler/action
            return handler(state, action as any);
        }
        return state;
    };
}

export const catalogReducer = createReducer<CatalogViewerState, CatalogSupportedActions>({
    AlbumsAndMediasLoadedAction: reduceAlbumsAndMediasLoaded,
    AlbumsLoadedAction: reduceAlbumsLoaded,
    MediaFailedToLoadAction: reduceMediaFailedToLoad,
});