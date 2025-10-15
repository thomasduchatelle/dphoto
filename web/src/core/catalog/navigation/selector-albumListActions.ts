import {CatalogViewerState, albumIsOwnedByCurrentUser} from "../language";
import {AlbumListActionsProps} from "../../../components/albums/AlbumsListActions";
import {displayedAlbumSelector} from "../language/selector-displayedAlbum";

export function albumListActionsSelector(state: CatalogViewerState): AlbumListActionsProps {
    const deleteButtonEnabled = state.allAlbums.some(album => albumIsOwnedByCurrentUser(album));
    const {displayedAlbumIdIsOwned} = displayedAlbumSelector(state);

    return {
        selected: state.albumFilter,
        options: state.albumFilterOptions,
        displayedAlbumIdIsOwned,
        deleteButtonEnabled,
    };
}
