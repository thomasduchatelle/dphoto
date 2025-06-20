import {CatalogViewerState, Album, Media, AlbumId} from "../language";
import {refreshFilters} from "../common/utils";

export interface AlbumDatesUpdated {
    type: "AlbumDatesUpdated";
    albums: Album[];
    medias: Media[];
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

    const mediasWithinDays = medias.reduce((acc, media) => {
        const dayKey = media.time.toDateString();
        const existingDay = acc.find(day => day.day.toDateString() === dayKey);
        
        if (existingDay) {
            existingDay.medias.push(media);
        } else {
            acc.push({
                day: new Date(media.time.getFullYear(), media.time.getMonth(), media.time.getDate()),
                medias: [media]
            });
        }
        
        return acc;
    }, [] as {day: Date, medias: Media[]}[]);

    return {
        ...current,
        allAlbums: albums,
        albumFilterOptions,
        albumFilter,
        albums: filteredAlbums,
        medias: mediasWithinDays,
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
