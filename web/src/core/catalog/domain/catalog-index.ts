import { AlbumsAndMediasLoadedAction, albumsAndMediasLoadedAction } from "./action-albumsAndMediasLoadedAction";
import { AlbumsLoadedAction, albumsLoadedAction } from "./action-albumsLoadedAction";
import { MediaFailedToLoadAction, mediaFailedToLoadAction } from "./action-mediaFailedToLoadAction";
import { NoAlbumAvailableAction, noAlbumAvailableAction } from "./action-noalbumavailableaction";
import { StartLoadingMediasAction, startLoadingMediasAction } from "./action-startloadingmediasaction";
import { catalogReducer } from "./catalog-reducer-v2";

export type CatalogSupportedActions =
    | AlbumsAndMediasLoadedAction
    | AlbumsLoadedAction
    | MediaFailedToLoadAction
    | NoAlbumAvailableAction
    | StartLoadingMediasAction;

export const catalogActions = {
    albumsAndMediasLoadedAction,
    albumsLoadedAction,
    mediaFailedToLoadAction,
    noAlbumAvailableAction,
    startLoadingMediasAction,
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
};
