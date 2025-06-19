import {AlbumId} from "../../language";
import {editAlbumDatesDialogOpened, EditAlbumDatesDialogOpened} from "./action-editAlbumDatesDialogOpened";
import {CatalogViewerState} from "../language";
import {ThunkDeclaration} from "../../thunk-engine";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {albumIdEquals} from "../language/utils-albumIdEquals";

export function openEditAlbumDatesDialogThunk(
    dispatch: (action: EditAlbumDatesDialogOpened) => void,
    albumId: AlbumId,
    state: CatalogViewerState, // Pass state to check for album existence
): void {
    // Only open the dialog if the album actually exists in the current state
    const albumExists = state.albums.some(album => albumIdEquals(album.albumId, albumId));
    if (albumExists) {
        dispatch(editAlbumDatesDialogOpened(albumId));
    }
}

export const openEditAlbumDatesDialogDeclaration: ThunkDeclaration<
    CatalogViewerState,
    {}, // No specific partial state needed from selector for this thunk
    (albumId: AlbumId) => void,
    CatalogFactoryArgs
> = {
    selector: (state: CatalogViewerState) => ({}), // No specific data needed from state for the factory
    factory: ({dispatch, app, state}) => {
        // Bind the thunk with dispatch and the current state
        return (albumId: AlbumId) => openEditAlbumDatesDialogThunk(dispatch, albumId, state);
    },
};
