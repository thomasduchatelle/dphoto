import type {ThunkDeclaration} from "../../thunk-engine";
import type {CatalogViewerAction, CatalogViewerState} from "../domain";
import {catalogActions} from "../domain";
import type {CatalogFactoryArgs} from "./catalog-factory-args";

export function openDeleteAlbumDialogThunk(
    dispatch: (action: CatalogViewerAction) => void
): void {
    console.debug("Opening delete album dialog");
    dispatch(catalogActions.openDeleteAlbumDialogAction());
}

export const openDeleteAlbumDialogDeclaration: ThunkDeclaration<
    CatalogViewerState,
    {},
    () => void,
    CatalogFactoryArgs
> = {
    selector: (_state: CatalogViewerState) => ({}),
    factory: ({dispatch}) => openDeleteAlbumDialogThunk.bind(null, dispatch),
};
