import {AlbumsAndMediasLoadedAction, albumsAndMediasLoadedAction, reduceAlbumsAndMediasLoaded} from "./catalog-action-AlbumsAndMediasLoadedAction";
import {AlbumsLoadedAction, albumsLoadedAction, reduceAlbumsLoaded} from "./catalog-action-albumsLoadedAction";
import {createReducer} from "./catalog-reducer-v2";
import {CatalogViewerState} from "./catalog-state";

export type CatalogSupportedActions =
    | AlbumsAndMediasLoadedAction
    | AlbumsLoadedAction;

export const catalogActions = {
    albumsAndMediasLoadedAction,
    albumsLoadedAction,
};

export const catalogReducer = createReducer<CatalogViewerState, CatalogSupportedActions>({
    AlbumsAndMediasLoadedAction: reduceAlbumsAndMediasLoaded,
    AlbumsLoadedAction: reduceAlbumsLoaded,
});

export type {
    AlbumsAndMediasLoadedAction,
    AlbumsLoadedAction,
};
