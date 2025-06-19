import {AlbumId} from "../../language";
import {editAlbumDatesDialogOpened, EditAlbumDatesDialogOpened} from "./action-editAlbumDatesDialogOpened";
import {CatalogViewerState} from "../language";
import {ThunkDeclaration} from "../../thunk-engine";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";

export function openEditAlbumDatesDialogThunk(
    dispatch: (action: EditAlbumDatesDialogOpened) => void,
    albumId: AlbumId,
): void {
    dispatch(editAlbumDatesDialogOpened(albumId));
}

export const openEditAlbumDatesDialogDeclaration: ThunkDeclaration<
    CatalogViewerState,
    {},
    (albumId: AlbumId) => void,
    CatalogFactoryArgs
> = {
    selector: (state: CatalogViewerState) => ({}),
    factory: ({dispatch, app}) => {
        return (albumId: AlbumId) => openEditAlbumDatesDialogThunk(dispatch, albumId);
    },
};
