import {AlbumsAndMediasLoadedAction, albumsAndMediasLoadedAction} from "./catalog-action-AlbumsAndMediasLoadedAction";
import {AlbumsLoadedAction, albumsLoadedAction} from "./catalog-action-albumsLoadedAction";
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
