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
import {CatalogViewerState} from "./language";
import {
    EditAlbumDatesDialogClosed,
    editAlbumDatesDialogClosedReducerRegistration
} from "./album-edit-dates/action-editAlbumDatesDialogClosed"; // New import
import {
    EditAlbumDatesDialogOpened,
    editAlbumDatesDialogOpenedReducerRegistration
} from "./album-edit-dates/action-editAlbumDatesDialogOpened"; // New import

export * from "./album-delete/selector-deleteDialogSelector";
export * from "./sharing/selector-sharingDialogSelector";
export * from "./navigation";
export * from "./album-edit-dates/selector-editAlbumDatesDialogSelector"; // New export

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
    EditAlbumDatesDialogOpened, // New type export
    EditAlbumDatesDialogClosed, // New type export
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
    | EditAlbumDatesDialogOpened // New action type
    | EditAlbumDatesDialogClosed // New action type

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
    editAlbumDatesDialogOpenedReducerRegistration, // New reducer registration
    editAlbumDatesDialogClosedReducerRegistration, // New reducer registration
];

function buildHandlers() {
    const handlers: any = {};
    for (const register of reducerRegistrations) {
        register(handlers);
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
