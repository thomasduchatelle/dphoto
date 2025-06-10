import type {ThunkDeclaration} from "../../thunk-engine";
import type {CatalogViewerState} from "../language";
import {DeleteAlbumDialogClosed, deleteAlbumDialogClosed} from "./action-deleteAlbumDialogClosed";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";

export function closeDeleteAlbumDialogThunk(
    dispatch: (action: DeleteAlbumDialogClosed) => void
): void {
    dispatch(deleteAlbumDialogClosed());
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
