import {AlbumFilterCriterion, AlbumFilterEntry} from "../language";

export interface AlbumListActionsProps {
    albumFilter: AlbumFilterEntry;
    albumFilterOptions: AlbumFilterEntry[];
    displayedAlbumIdIsOwned: boolean;
    hasAlbumsToDelete: boolean;
    canCreateAlbums: boolean;
}
