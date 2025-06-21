import {Album, CatalogViewerState, MediaWithinADay} from "../language";
import {refreshFilters} from "../common/utils";

export interface AlbumDatesUpdated {
    type: "AlbumDatesUpdated";
    albums: Album[];
    medias: MediaWithinADay[];
}

export function albumDatesUpdated(props: Omit<AlbumDatesUpdated, "type">): AlbumDatesUpdated {
    return {
        ...props,
        type: "AlbumDatesUpdated",
    };
}

export function reduceAlbumDatesUpdated(
    current: CatalogViewerState,
    {albums, medias}: AlbumDatesUpdated,
): CatalogViewerState {
    const {albumFilterOptions, albumFilter, albums: filteredAlbums} = refreshFilters(
        current.currentUser,
        current.albumFilter,
        albums
    );

    return {
        ...current,
        allAlbums: albums,
        albumFilterOptions,
        albumFilter,
        albums: filteredAlbums,
        medias: medias,
        mediasLoadedFromAlbumId: current.editDatesDialog?.albumId, // Use albumId from dialog state
        albumsLoaded: true,
        mediasLoaded: true,
        editDatesDialog: undefined,
    };
}

export function albumDatesUpdatedReducerRegistration(handlers: any) {
    handlers["AlbumDatesUpdated"] = reduceAlbumDatesUpdated as (
        state: CatalogViewerState,
        action: AlbumDatesUpdated
    ) => CatalogViewerState;
}
