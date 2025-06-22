import {
    AlbumAccessGranted,
    albumAccessGrantedReducerRegistration,
    AlbumAccessRevoked,
    albumAccessRevokedReducerRegistration,
    SharingModalClosed,
    sharingModalClosedReducerRegistration,
    SharingModalErrorOccurred,
    sharingModalErrorOccurredReducerRegistration,
    SharingModalOpened,
    sharingModalOpenedReducerRegistration,
} from "./sharing";
import {
    AlbumsAndMediasLoaded,
    albumsAndMediasLoadedReducerRegistration,
    AlbumsFiltered,
    albumsFilteredReducerRegistration,
    AlbumsLoaded,
    albumsLoadedReducerRegistration,
    MediaLoadFailed,
    mediaLoadFailedReducerRegistration,
    MediasLoaded,
    mediasLoadedReducerRegistration,
    MediasLoadingStarted,
    mediasLoadingStartedReducerRegistration,
    NoAlbumAvailable,
    noAlbumAvailableReducerRegistration
} from "./navigation";
import {
    AlbumDeleted,
    albumDeletedReducerRegistration,
    AlbumDeleteFailed,
    albumDeleteFailedReducerRegistration,
    DeleteAlbumDialogClosed,
    deleteAlbumDialogClosedReducerRegistration,
    DeleteAlbumDialogOpened,
    deleteAlbumDialogOpenedReducerRegistration,
    DeleteAlbumStarted,
    deleteAlbumStartedReducerRegistration
} from "./album-delete";
import {EditDatesDialogClosed, EditDatesDialogOpened} from "./album-edit-dates";
import {CatalogViewerState} from "./language";
import {ActionWithReducer} from "./common/action-factory";

export * from "./album-delete/selector-deleteDialogSelector";
export * from "./album-edit-dates/selector-editDatesDialogSelector";
export * from "./sharing/selector-sharingDialogSelector";
export * from "./navigation";

export type {
    AlbumAccessGranted,
    AlbumsAndMediasLoaded,
    AlbumsFiltered,
    AlbumsLoaded,
    AlbumDeleted,
    SharingModalClosed,
    MediaLoadFailed,
    MediasLoaded,
    NoAlbumAvailable,
    SharingModalOpened,
    AlbumAccessRevoked,
    SharingModalErrorOccurred,
    MediasLoadingStarted,
    DeleteAlbumDialogOpened,
    AlbumDeleteFailed,
    DeleteAlbumDialogClosed,
    DeleteAlbumStarted,
    EditDatesDialogOpened,
    EditDatesDialogClosed,
};

// Legacy action types for backward compatibility
export type CatalogViewerAction =
    ActionWithReducer<CatalogViewerState>
    | AlbumAccessGranted
    | AlbumsAndMediasLoaded
    | AlbumsFiltered
    | AlbumsLoaded
    | AlbumDeleted
    | SharingModalClosed
    | MediaLoadFailed
    | MediasLoaded
    | NoAlbumAvailable
    | SharingModalOpened
    | AlbumAccessRevoked
    | SharingModalErrorOccurred
    | MediasLoadingStarted
    | DeleteAlbumDialogOpened
    | AlbumDeleteFailed
    | DeleteAlbumDialogClosed
    | DeleteAlbumStarted

const reducerRegistrations = [
    albumAccessGrantedReducerRegistration,
    albumsAndMediasLoadedReducerRegistration,
    albumsLoadedReducerRegistration,
    albumDeletedReducerRegistration,
    mediaLoadFailedReducerRegistration,
    noAlbumAvailableReducerRegistration,
    mediasLoadingStartedReducerRegistration,
    albumsFilteredReducerRegistration,
    sharingModalOpenedReducerRegistration,
    albumAccessRevokedReducerRegistration,
    sharingModalClosedReducerRegistration,
    sharingModalErrorOccurredReducerRegistration,
    mediasLoadedReducerRegistration,
    deleteAlbumDialogOpenedReducerRegistration,
    albumDeleteFailedReducerRegistration,
    deleteAlbumDialogClosedReducerRegistration,
    deleteAlbumStartedReducerRegistration,
];

function buildHandlers() {
    const handlers: any = {};

    // Register old-style reducers
    for (const register of reducerRegistrations) {
        register(handlers);
    }

    return handlers;
}

function createGenericReducer<TState>(
    legacyHandlers: Record<string, (state: TState, action: any) => TState>
): (state: TState, action: ActionWithReducer<TState> | { type: string }) => TState {
    return (state: TState, action: ActionWithReducer<TState> | { type: string }): TState => {
        // Check if action has a built-in reducer
        if ('reducer' in action && typeof action.reducer === 'function') {
            return action.reducer(state, action as ActionWithReducer<TState>);
        }

        // Fall back to legacy handlers
        const handler = legacyHandlers[action.type];
        if (handler) {
            return handler(state, action);
        }

        return state;
    };
}

function createCatalogReducer(): (state: CatalogViewerState, action: ActionWithReducer<CatalogViewerState> | CatalogViewerAction) => CatalogViewerState {
    const handlers = buildHandlers();
    return createGenericReducer(handlers);
}

export const catalogReducer = createCatalogReducer();
