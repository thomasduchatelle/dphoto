import {AlbumsAndMediasLoadedAction, albumsAndMediasLoadedAction} from "./action-albumsAndMediasLoadedAction";
import {AlbumsLoadedAction, albumsLoadedAction} from "./catalog-action-AlbumsLoadedAction";
import {MediaFailedToLoadAction, mediaFailedToLoadAction} from "./catalog-action-MediaFailedToLoadAction";

export type CatalogSupportedActions =
    | AlbumsAndMediasLoadedAction
    | AlbumsLoadedAction
    | MediaFailedToLoadAction;

export const catalogActions = {
    albumsAndMediasLoadedAction,
    albumsLoadedAction,
    mediaFailedToLoadAction,
};

export type {
    AlbumsAndMediasLoadedAction,
    AlbumsLoadedAction,
    MediaFailedToLoadAction,
};
