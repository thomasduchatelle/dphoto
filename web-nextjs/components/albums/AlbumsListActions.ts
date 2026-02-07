import {AlbumFilterEntry} from '@/domains/catalog/language/catalog-state';

export interface AlbumListActionsProps {
    albumFilter: AlbumFilterEntry;
    albumFilterOptions: AlbumFilterEntry[];
    displayedAlbumIdIsOwned: boolean;
    hasAlbumsToDelete: boolean;
    canCreateAlbums: boolean;
}
