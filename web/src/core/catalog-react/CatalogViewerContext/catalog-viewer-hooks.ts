import {CatalogViewerState} from "./catalog-viewer-state";
import {CatalogViewerContext} from "./catalog-viewer-context";
import {useContext} from "react";

export const useCatalogViewerState = (): CatalogViewerState => {
    const {state} = useContext(CatalogViewerContext)
    return state
}