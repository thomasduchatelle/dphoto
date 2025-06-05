import {AddSharingAction, addSharingAction, addSharingReducerRegistration} from "./sharing-addSharingAction";
import {AlbumsAndMediasLoadedAction, albumsAndMediasLoadedAction, albumsAndMediasLoadedReducerRegistration} from "./content-albumsAndMediasLoadedAction";
import {AlbumsLoadedAction, albumsLoadedAction, albumsLoadedReducerRegistration} from "./content-albumsLoadedAction";
import {AlbumDeletedAction, albumDeletedAction, albumDeletedReducerRegistration} from "./delete-albumDeletedAction";
import {MediaFailedToLoadAction, mediaFailedToLoadAction, mediaFailedToLoadReducerRegistration} from "./content-mediaFailedToLoadAction";
import {NoAlbumAvailableAction, noAlbumAvailableAction, noAlbumAvailableReducerRegistration} from "./content-noAlbumAvailableAction";
import {StartLoadingMediasAction, startLoadingMediasAction, startLoadingMediasReducerRegistration} from "./content-startLoadingMediasAction";
import {AlbumsFilteredAction, albumsFilteredAction, albumsFilteredReducerRegistration} from "./content-albumsFilteredAction";
import {OpenSharingModalAction, openSharingModalAction, openSharingModalReducerRegistration} from "./sharing-openSharingModalAction";
import {RemoveSharingAction, removeSharingAction, removeSharingReducerRegistration} from "./sharing-removeSharingAction";
import {CloseSharingModalAction, closeSharingModalAction, closeSharingModalReducerRegistration} from "./sharing-closeSharingModalAction";
import {SharingModalErrorAction, sharingModalErrorAction, sharingModalErrorReducerRegistration} from "./sharing-sharingModalErrorAction";
import {MediasLoadedAction, mediasLoadedAction, mediasLoadedReducerRegistration} from "./content-mediasLoadedAction";
import {CatalogViewerState} from "../catalog-state";
import {OpenDeleteAlbumDialogAction, openDeleteAlbumDialogAction, openDeleteAlbumDialogReducerRegistration} from "./delete-openDeleteAlbumDialog";
import {AlbumFailedToDeleteAction, albumFailedToDeleteAction, albumFailedToDeleteReducerRegistration} from "./delete-albumFailedToDeleteAction";
import {CloseDeleteAlbumDialogAction, closeDeleteAlbumDialogAction, closeDeleteAlbumDialogReducerRegistration} from "./delete-closeDeleteAlbumDialog";
import {StartDeleteAlbumAction, startDeleteAlbumAction, startDeleteAlbumReducerRegistration} from "./delete-startDeleteAlbum";

export * from "./selector-deleteDialogSelector";
export * from "./selector-sharingDialogSelector";
export * from "./catalog-common-modifiers";

export const catalogActions = {
    openSharingModalAction,
    addSharingAction,
    removeSharingAction,
    closeSharingModalAction,
    sharingModalErrorAction,
    albumsAndMediasLoadedAction,
    albumsLoadedAction,
    albumDeletedAction,
    mediaFailedToLoadAction,
    noAlbumAvailableAction,
    startLoadingMediasAction,
    albumsFilteredAction,
    mediasLoadedAction,
    openDeleteAlbumDialogAction,
    albumFailedToDeleteAction,
    closeDeleteAlbumDialogAction,
    startDeleteAlbumAction,
};

export type {
    AddSharingAction,
    AlbumsAndMediasLoadedAction,
    AlbumsFilteredAction,
    AlbumsLoadedAction,
    AlbumDeletedAction,
    CloseSharingModalAction,
    MediaFailedToLoadAction,
    MediasLoadedAction,
    NoAlbumAvailableAction,
    OpenSharingModalAction,
    RemoveSharingAction,
    SharingModalErrorAction,
    StartLoadingMediasAction,
    OpenDeleteAlbumDialogAction,
    AlbumFailedToDeleteAction,
    CloseDeleteAlbumDialogAction,
    StartDeleteAlbumAction,
};

export type CatalogViewerAction =
    AddSharingAction
    | AlbumsAndMediasLoadedAction
    | AlbumsFilteredAction
    | AlbumsLoadedAction
    | AlbumDeletedAction
    | CloseSharingModalAction
    | MediaFailedToLoadAction
    | MediasLoadedAction
    | NoAlbumAvailableAction
    | OpenSharingModalAction
    | RemoveSharingAction
    | SharingModalErrorAction
    | StartLoadingMediasAction
    | OpenDeleteAlbumDialogAction
    | AlbumFailedToDeleteAction
    | CloseDeleteAlbumDialogAction
    | StartDeleteAlbumAction

const reducerRegistrations = [
    addSharingReducerRegistration,
    albumsAndMediasLoadedReducerRegistration,
    albumsLoadedReducerRegistration,
    albumDeletedReducerRegistration,
    mediaFailedToLoadReducerRegistration,
    noAlbumAvailableReducerRegistration,
    startLoadingMediasReducerRegistration,
    albumsFilteredReducerRegistration,
    openSharingModalReducerRegistration,
    removeSharingReducerRegistration,
    closeSharingModalReducerRegistration,
    sharingModalErrorReducerRegistration,
    mediasLoadedReducerRegistration,
    openDeleteAlbumDialogReducerRegistration,
    albumFailedToDeleteReducerRegistration,
    closeDeleteAlbumDialogReducerRegistration,
    startDeleteAlbumReducerRegistration,
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
