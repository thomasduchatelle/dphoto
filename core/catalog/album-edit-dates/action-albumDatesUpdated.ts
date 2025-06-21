export interface AlbumDatesUpdated {
    type: "AlbumDatesUpdated";
    albums: Album[];
    medias: Media[];
    updatedAlbumId: AlbumId;
}

export function albumDatesUpdated(props: Omit<AlbumDatesUpdated, "type">): AlbumDatesUpdated {
    return {
        ...props,
        type: "AlbumDatesUpdated",
    };
}

export function reduceAlbumDatesUpdated(
    current: CatalogViewerState,
    {albums, medias, updatedAlbumId}: AlbumDatesUpdated,
): CatalogViewerState {
    return {
        ...current,
        albums,
        allAlbums: albums,
        medias,
        mediasLoaded: true,
        mediasLoadedFromAlbumId: updatedAlbumId,
        editDatesDialog: undefined,
    };
}

export function albumDatesUpdatedReducerRegistration(handlers: any) {
    handlers["AlbumDatesUpdated"] = reduceAlbumDatesUpdated as (
        state: CatalogViewerState,
        action: AlbumDatesUpdated
    ) => CatalogViewerState;
}
