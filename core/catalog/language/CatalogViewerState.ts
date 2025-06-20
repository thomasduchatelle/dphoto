export interface CatalogViewerState {
    albums: Album[];
    allAlbums: Album[];
    medias: Media[];
    mediasLoaded: boolean;
    mediasLoadedFromAlbumId?: AlbumId;
    albumsLoaded: boolean;
    error?: string;
    createDialog?: CreateDialogState;
    deleteDialog?: DeleteDialogState;
    shareModal?: ShareModalState;
    editDatesDialog?: EditDatesDialogState;
}

export interface EditDatesDialogState {
    albumId: AlbumId;
    albumName: string;
    startDate: Date;
    endDate: Date;
    isLoading: boolean;
}
