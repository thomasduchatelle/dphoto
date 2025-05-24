import {AddSharingAction, addSharingAction, addSharingReducerRegistration} from "./action-addSharingAction";
import {AlbumsAndMediasLoadedAction, albumsAndMediasLoadedAction, albumsAndMediasLoadedReducerRegistration} from "./action-albumsAndMediasLoadedAction";
import {AlbumsLoadedAction, albumsLoadedAction, albumsLoadedReducerRegistration} from "./action-albumsLoadedAction";
import {MediaFailedToLoadAction, mediaFailedToLoadAction, mediaFailedToLoadReducerRegistration} from "./action-mediaFailedToLoadAction";
import {NoAlbumAvailableAction, noAlbumAvailableAction, noAlbumAvailableReducerRegistration} from "./action-noAlbumAvailableAction";
import {StartLoadingMediasAction, startLoadingMediasAction, startLoadingMediasReducerRegistration} from "./action-startLoadingMediasAction";
import {AlbumsFilteredAction, albumsFilteredAction, albumsFilteredReducerRegistration} from "./action-albumsFilteredAction";
import {OpenSharingModalAction, openSharingModalAction, openSharingModalReducerRegistration} from "./action-openSharingModalAction";
import {RemoveSharingAction, removeSharingAction, removeSharingReducerRegistration} from "./action-removeSharingAction";
import {CloseSharingModalAction, closeSharingModalAction, closeSharingModalReducerRegistration} from "./action-closeSharingModalAction";
import {SharingModalErrorAction, sharingModalErrorAction, sharingModalErrorReducerRegistration} from "./action-sharingModalErrorAction";
import {MediasLoadedAction, mediasLoadedAction, mediasLoadedReducerRegistration} from "./action-mediasLoadedAction";
import {CatalogViewerState} from "./catalog-state";

export const catalogActions = {
    openSharingModalAction,
    addSharingAction,
    removeSharingAction,
    closeSharingModalAction,
    sharingModalErrorAction,
    albumsAndMediasLoadedAction,
    albumsLoadedAction,
    mediaFailedToLoadAction,
    noAlbumAvailableAction,
    startLoadingMediasAction,
    albumsFilteredAction,
    mediasLoadedAction,
};

export type {
    AddSharingAction,
    AlbumsAndMediasLoadedAction,
    AlbumsFilteredAction,
    AlbumsLoadedAction,
    CloseSharingModalAction,
    MediaFailedToLoadAction,
    MediasLoadedAction,
    NoAlbumAvailableAction,
    OpenSharingModalAction,
    RemoveSharingAction,
    SharingModalErrorAction,
    StartLoadingMediasAction,
};

export type CatalogViewerAction =
    AddSharingAction
    | AlbumsAndMediasLoadedAction
    | AlbumsFilteredAction
    | AlbumsLoadedAction
    | CloseSharingModalAction
    | MediaFailedToLoadAction
    | MediasLoadedAction
    | NoAlbumAvailableAction
    | OpenSharingModalAction
    | RemoveSharingAction
    | SharingModalErrorAction
    | StartLoadingMediasAction

const reducerRegistrations = [
    addSharingReducerRegistration,
    albumsAndMediasLoadedReducerRegistration,
    albumsLoadedReducerRegistration,
    mediaFailedToLoadReducerRegistration,
    noAlbumAvailableReducerRegistration,
    startLoadingMediasReducerRegistration,
    albumsFilteredReducerRegistration,
    openSharingModalReducerRegistration,
    removeSharingReducerRegistration,
    closeSharingModalReducerRegistration,
    sharingModalErrorReducerRegistration,
    mediasLoadedReducerRegistration,
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
