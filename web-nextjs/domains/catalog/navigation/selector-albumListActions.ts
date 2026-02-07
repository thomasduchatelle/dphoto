import {AlbumFilterEntry, albumIsOwnedByCurrentUser, CatalogViewerState} from "../language";
import {displayedAlbumSelector} from "../language/selector-displayedAlbum";

export interface AlbumListActionsProps {
    albumFilter: AlbumFilterEntry;
    albumFilterOptions: AlbumFilterEntry[];
    displayedAlbumIdIsOwned: boolean;
    hasAlbumsToDelete: boolean;
    canCreateAlbums: boolean;
}

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
