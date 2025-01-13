import {CatalogViewerContext} from "./CatalogViewerProvider";
import {useContext} from "react";
import {AlbumId, CatalogViewerState} from "../../catalog";
import {CatalogHandlers} from "./CatalogViewerStateWithDispatch";

export const useCatalogContext = (): { state: CatalogViewerState, handlers: CatalogHandlers, selectedAlbumId?: AlbumId } => {
    const {state, handlers, selectedAlbumId} = useContext(CatalogViewerContext)
    return {state, handlers, selectedAlbumId}
}