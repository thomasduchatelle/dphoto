import {albumIsOwnedByCurrentUser, CatalogViewerState} from "../language";
import {AlbumListActionsProps} from "./types-albumListActions";
import {displayedAlbumSelector} from "../language/selector-displayedAlbum";

export function albumListActionsSelector(state: CatalogViewerState): AlbumListActionsProps {
    const deleteButtonEnabled = state.allAlbums.some(album => albumIsOwnedByCurrentUser(album));
    const {displayedAlbumIdIsOwned} = displayedAlbumSelector(state);
    const createButtonEnabled = state.currentUser.isOwner;

    return {
        albumFilter: state.albumFilter,
        albumFilterOptions: state.albumFilterOptions,
        displayedAlbumIdIsOwned,
        hasAlbumsToDelete: deleteButtonEnabled,
        canCreateAlbums: createButtonEnabled,
    };
}
