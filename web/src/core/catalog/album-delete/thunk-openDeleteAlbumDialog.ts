import type {CatalogViewerState} from "../language";
import {deleteAlbumDialogOpened} from "./action-deleteAlbumDialogOpened";
import type {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {ThunkDeclaration} from "src/libs/dthunks";

export function openDeleteAlbumDialogThunk(
    dispatch: (action: ReturnType<typeof deleteAlbumDialogOpened>) => void
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
