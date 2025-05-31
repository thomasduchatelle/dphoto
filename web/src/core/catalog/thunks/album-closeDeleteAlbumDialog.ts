import type {ThunkDeclaration} from "../../thunk-engine";
import type {CatalogViewerAction, CatalogViewerState} from "../domain";
import {catalogActions} from "../domain";
import {CatalogFactoryArgs} from "./catalog-factory-args";

export function closeDeleteAlbumDialogThunk(
    dispatch: (action: CatalogViewerAction) => void
): void {
    dispatch(catalogActions.closeDeleteAlbumDialogAction());
}

export const closeDeleteAlbumDialogDeclaration: ThunkDeclaration<
    CatalogViewerState,
    {},
    () => void,
    CatalogFactoryArgs
> = {
    selector: (_state: CatalogViewerState) => ({}),
    factory: ({dispatch}) => closeDeleteAlbumDialogThunk.bind(null, dispatch),
};
