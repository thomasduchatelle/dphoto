import { AlbumsAndMediasLoadedAction, albumsAndMediasLoadedAction, reduceAlbumsAndMediasLoaded } from "./catalog-action-AlbumsAndMediasLoadedAction";
import { AlbumsLoadedAction, albumsLoadedAction, reduceAlbumsLoaded } from "./catalog-action-albumsLoadedAction";
import { MediaFailedToLoadAction, mediaFailedToLoadAction, reduceMediaFailedToLoad } from "./catalog-action-MediaFailedToLoadAction";
import { createReducer } from "./catalog-reducer-v2";
import { CatalogViewerState } from "./catalog-state";

export type CatalogSupportedActions =
    | AlbumsAndMediasLoadedAction
    | AlbumsLoadedAction
    | MediaFailedToLoadAction;

export const catalogActions = {
    albumsAndMediasLoadedAction,
    albumsLoadedAction,
    mediaFailedToLoadAction,
};

export const catalogReducer = createReducer<CatalogViewerState, CatalogSupportedActions>({
    AlbumsAndMediasLoadedAction: reduceAlbumsAndMediasLoaded,
    AlbumsLoadedAction: reduceAlbumsLoaded,
    MediaFailedToLoadAction: reduceMediaFailedToLoad,
});

export type {
    AlbumsAndMediasLoadedAction,
    AlbumsLoadedAction,
    MediaFailedToLoadAction,
};
