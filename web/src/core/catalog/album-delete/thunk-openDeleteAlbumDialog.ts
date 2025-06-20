import type {ThunkDeclaration} from "../../thunk-engine";
import type {CatalogViewerState} from "../language";
import {DeleteAlbumDialogOpened, deleteAlbumDialogOpened} from "./action-deleteAlbumDialogOpened";
import type {CatalogFactoryArgs} from "../common/catalog-factory-args";

export function openDeleteAlbumDialogThunk(
    dispatch: (action: DeleteAlbumDialogOpened) => void
): void {
    console.debug("Opening delete album dialog");
    dispatch(deleteAlbumDialogOpened());
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
