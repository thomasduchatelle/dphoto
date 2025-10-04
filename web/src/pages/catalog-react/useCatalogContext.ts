import {CatalogViewerContext} from "./CatalogViewerProvider";
import {useContext} from "react";
import {AlbumId, CatalogThunksInterface, CatalogViewerState} from "../../core/catalog";

export const useCatalogContext = (): {
    state: CatalogViewerState,
    handlers: Omit<CatalogThunksInterface, "onPageRefresh">,
    selectedAlbumId?: AlbumId
} => {
    const {state, handlers, selectedAlbumId} = useContext(CatalogViewerContext)
    if (!handlers) {
        throw new Error("CatalogViewerContext not initialized. Ensure CatalogViewerProvider is used in the component tree.");
    }
    return {state, handlers, selectedAlbumId};
}