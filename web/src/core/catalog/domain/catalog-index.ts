import { AlbumsAndMediasLoadedAction, albumsAndMediasLoadedAction } from "./action-albumsAndMediasLoadedAction";
import { AlbumsLoadedAction, albumsLoadedAction } from "./action-albumsLoadedAction";
import { MediaFailedToLoadAction, mediaFailedToLoadAction } from "./action-mediaFailedToLoadAction";
import { NoAlbumAvailableAction, noAlbumAvailableAction } from "./action-noalbumavailableaction";
import { StartLoadingMediasAction, startLoadingMediasAction } from "./action-startloadingmediasaction";
import { AlbumsFilteredAction, albumsFilteredAction } from "./action-albumsfilteredaction";
import { catalogReducer } from "./catalog-reducer-v2";

export type CatalogSupportedActions =
    | AlbumsAndMediasLoadedAction
    | AlbumsLoadedAction
    | MediaFailedToLoadAction
    | NoAlbumAvailableAction
    | StartLoadingMediasAction
    | AlbumsFilteredAction;

export const catalogActions = {
    albumsAndMediasLoadedAction,
    albumsLoadedAction,
    mediaFailedToLoadAction,
    noAlbumAvailableAction,
    startLoadingMediasAction,
    albumsFilteredAction,
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
};
