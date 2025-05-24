import {AlbumsAndMediasLoadedAction, albumsAndMediasLoadedAction} from "./action-albumsAndMediasLoadedAction";
import {AlbumsLoadedAction, albumsLoadedAction} from "./action-albumsLoadedAction";
import {MediaFailedToLoadAction, mediaFailedToLoadAction} from "./action-mediaFailedToLoadAction";
import {NoAlbumAvailableAction, noAlbumAvailableAction} from "./action-noAlbumAvailableAction";
import {StartLoadingMediasAction, startLoadingMediasAction} from "./action-startLoadingMediasAction";
import {AlbumsFilteredAction, albumsFilteredAction} from "./action-albumsFilteredAction";
import {OpenSharingModalAction, openSharingModalAction} from "./action-openSharingModalAction";
import {AddSharingAction, addSharingAction} from "./action-addSharingAction";
import {RemoveSharingAction, removeSharingAction} from "./action-removeSharingAction";
import {catalogReducer} from "./catalog-reducer-v2";

export type CatalogSupportedActions =
    | AlbumsAndMediasLoadedAction
    | AlbumsLoadedAction
    | MediaFailedToLoadAction
    | NoAlbumAvailableAction
    | StartLoadingMediasAction
    | AlbumsFilteredAction
    | OpenSharingModalAction
    | AddSharingAction
    | RemoveSharingAction;

export const catalogActions = {
    albumsAndMediasLoadedAction,
    albumsLoadedAction,
    mediaFailedToLoadAction,
    noAlbumAvailableAction,
    startLoadingMediasAction,
    albumsFilteredAction,
    openSharingModalAction,
    addSharingAction,
    removeSharingAction,
};

export {
    catalogReducer,
};

export type {
    AlbumsAndMediasLoadedAction,
    AlbumsLoadedAction,
    MediaFailedToLoadAction,
    NoAlbumAvailableAction,
    StartLoadingMediasAction,
    AlbumsFilteredAction,
    OpenSharingModalAction,
    AddSharingAction,
    RemoveSharingAction,
};
