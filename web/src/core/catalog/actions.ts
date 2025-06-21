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
import {
    AlbumDatesUpdated,
    albumDatesUpdatedReducerRegistration,
    AlbumDatesUpdateStarted,
    albumDatesUpdateStartedReducerRegistration,
    EditDatesDialogClosed,
    editDatesDialogClosedReducerRegistration,
    EditDatesDialogEndDateUpdated,
    editDatesDialogEndDateUpdatedReducerRegistration,
    EditDatesDialogOpened,
    editDatesDialogOpenedReducerRegistration,
    EditDatesDialogStartDateUpdated,
    editDatesDialogStartDateUpdatedReducerRegistration
} from "./album-edit-dates";
import {CatalogViewerState} from "./language";

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
    AlbumDatesUpdateStarted,
    AlbumDatesUpdated,
};

export type CatalogViewerAction =
    AlbumAccessGranted
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
    | EditDatesDialogOpened
    | EditDatesDialogClosed
    | AlbumDatesUpdateStarted
    | AlbumDatesUpdated
    | EditDatesDialogStartDateUpdated
    | EditDatesDialogEndDateUpdated

import {editDatesDialogClosed} from "./album-edit-dates/action-editDatesDialogClosed";
import {editDatesDialogStartDateUpdated} from "./album-edit-dates/action-editDatesDialogStartDateUpdated";

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
    editDatesDialogOpenedReducerRegistration,
    albumDatesUpdateStartedReducerRegistration,
    albumDatesUpdatedReducerRegistration,
    editDatesDialogEndDateUpdatedReducerRegistration,
];

// New-style action creators with built-in reducers
const newStyleActions = [
    editDatesDialogClosed,
    editDatesDialogStartDateUpdated,
];

function buildHandlers() {
    const handlers: any = {};
    
    // Register old-style reducers
    for (const register of reducerRegistrations) {
        register(handlers);
    }
    
    // Register new-style action reducers
    for (const actionCreator of newStyleActions) {
        handlers[actionCreator.type] = actionCreator.reducer;
    }
    
    return handlers;
}

function createReducer<TState, TActions extends { type: string }>(
    handlers: {
        [K in TActions["type"]]: (state: TState, action: Extract<TActions, { type: K }>) => TState
    }
): (state: TState, action: TActions) => TState {
    return (state: TState, action: TActions): TState => {
        const handler = handlers[action.type as keyof typeof handlers];
        if (handler) {
            return handler(state, action as any);
        }
        return state;
    };
}

function createCatalogReducer(): (state: CatalogViewerState, action: CatalogViewerAction) => CatalogViewerState {
    const handlers = buildHandlers();
    return createReducer(handlers);
}

export const catalogReducer = createCatalogReducer();
