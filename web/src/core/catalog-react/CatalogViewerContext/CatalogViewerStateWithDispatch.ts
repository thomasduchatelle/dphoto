import {AlbumId, CatalogViewerAction, CatalogViewerState} from "../../catalog";
import {Dispatch} from "react";
import {CatalogHandlers} from "./CatalogViewerProvider";

export interface CatalogViewerStateWithDispatch {
    state: CatalogViewerState
    selectedAlbumId?: AlbumId // state managed from the URL
    dispatch: Dispatch<CatalogViewerAction>
    handlers: CatalogHandlers
}