import {AlbumFilterCriterion, AlbumId, CatalogViewerState, CreateAlbumListener} from "../../catalog";

export interface CatalogHandlers extends CreateAlbumListener {
    onAlbumFilterChange: (criterion: AlbumFilterCriterion) => void
}

export interface CatalogViewerStateWithDispatch {
    state: CatalogViewerState
    selectedAlbumId?: AlbumId // state managed from the URL
    handlers: CatalogHandlers
}