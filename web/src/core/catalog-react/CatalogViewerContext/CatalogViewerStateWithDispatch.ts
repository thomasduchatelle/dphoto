import {AlbumFilterCriterion, AlbumId, CatalogViewerAction, CatalogViewerState, CreateAlbumListener} from "../../catalog";
import {Dispatch} from "react";

export interface CatalogHandlers extends CreateAlbumListener {
    onAlbumFilterChange: (criterion: AlbumFilterCriterion) => void
}

export interface CatalogViewerStateWithDispatch {
    state: CatalogViewerState
    selectedAlbumId?: AlbumId // state managed from the URL
    dispatch: Dispatch<CatalogViewerAction> // TODO the dispatch should not be exposed ; only readonly state and handlers
    handlers: CatalogHandlers
}