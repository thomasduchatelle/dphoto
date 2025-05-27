import {CatalogViewerContext} from "./CatalogViewerProvider";
import {useContext} from "react";
import {AlbumId, CatalogViewerState} from "../../catalog";
import {CatalogThunksInterface} from "../../catalog/thunks";

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