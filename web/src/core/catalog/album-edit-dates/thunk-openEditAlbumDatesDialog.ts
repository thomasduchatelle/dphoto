import {AlbumId, CatalogViewerState, currentAlbumIdSelector} from "../language";
import {editAlbumDatesDialogOpened, EditAlbumDatesDialogOpened} from "./action-editAlbumDatesDialogOpened";
import {ThunkDeclaration} from "../../thunk-engine";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";

export function openEditAlbumDatesDialogThunk(
    dispatch: (action: EditAlbumDatesDialogOpened) => void,
    albumId?: AlbumId,
): void {
    if (albumId) {
        dispatch(editAlbumDatesDialogOpened(albumId));
    }
}

export const openEditAlbumDatesDialogDeclaration: ThunkDeclaration<
    CatalogViewerState,
    { albumId?: AlbumId },
    () => void,
    CatalogFactoryArgs
> = {
    selector: (state: CatalogViewerState) => ({albumId: currentAlbumIdSelector(state)}),
    factory: ({dispatch, partialState: {albumId}}) => {
        return openEditAlbumDatesDialogThunk.bind(null, dispatch, albumId);
    },
};
