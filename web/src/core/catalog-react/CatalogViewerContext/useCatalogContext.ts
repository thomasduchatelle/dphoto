import {CatalogHandlers, CatalogViewerContext} from "./CatalogViewerProvider";
import {useContext} from "react";
import {AlbumId, CatalogViewerState} from "../../catalog";

export const useCatalogContext = (): { state: CatalogViewerState, handlers: CatalogHandlers, selectedAlbumId?: AlbumId } => {
    const {state, handlers, selectedAlbumId} = useContext(CatalogViewerContext)
    return {state, handlers, selectedAlbumId}
}