import { AlbumsAndMediasLoadedAction, albumsAndMediasLoadedAction } from "./action-albumsAndMediasLoadedAction";
import { AlbumsLoadedAction, albumsLoadedAction } from "./action-albumsLoadedAction";
import { MediaFailedToLoadAction, mediaFailedToLoadAction } from "./action-mediaFailedToLoadAction";
import { NoAlbumAvailableAction, noAlbumAvailableAction } from "./action-noalbumavailableaction";
import { catalogReducer } from "./catalog-reducer";

export type CatalogSupportedActions =
    | AlbumsAndMediasLoadedAction
    | AlbumsLoadedAction
    | MediaFailedToLoadAction
    | NoAlbumAvailableAction;

export const catalogActions = {
    albumsAndMediasLoadedAction,
    albumsLoadedAction,
    mediaFailedToLoadAction,
    noAlbumAvailableAction,
};

export {
    catalogReducer,
};

export type {
    AlbumsAndMediasLoadedAction,
    AlbumsLoadedAction,
    MediaFailedToLoadAction,
    NoAlbumAvailableAction,
};
