import {Album, AlbumId, albumIdEquals, CatalogViewerState, RedirectToAlbumIdAction} from "../language";
import {albumFilterAreCriterionEqual, ALL_ALBUMS_FILTER_CRITERION, DEFAULT_ALBUM_FILTER_ENTRY, refreshFilters} from "../navigation";
import {createAction} from "@light-state";

interface AlbumDeletedPayload extends RedirectToAlbumIdAction {
    albums: Album[];
    redirectTo?: AlbumId;
}

export const albumDeleted = createAction<CatalogViewerState, AlbumDeletedPayload>(
    "albumDeleted",
    (current: CatalogViewerState, {albums: newAlbums, redirectTo}: AlbumDeletedPayload) => {
        let {albumFilterOptions, albumFilter, albums} = refreshFilters(current.currentUser, current.albumFilter, newAlbums);

        if (
            redirectTo &&
            !newAlbums.some(album => albumIdEquals(album.albumId, redirectTo))
        ) {
            albumFilter =
                albumFilterOptions.find(option =>
                    albumFilterAreCriterionEqual(option.criterion, ALL_ALBUMS_FILTER_CRITERION)
                ) ?? DEFAULT_ALBUM_FILTER_ENTRY;
            albums = newAlbums
        }

        return {
            ...current,
            albumFilterOptions,
            albumFilter,
            allAlbums: newAlbums,
            albums: albums,
            error: undefined,
            albumsLoaded: true,
            deleteDialog: undefined,
        };
    }
);

export type AlbumDeleted = ReturnType<typeof albumDeleted>;
