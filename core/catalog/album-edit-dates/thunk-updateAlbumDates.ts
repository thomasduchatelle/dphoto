export interface UpdateAlbumDatesPort {
    updateAlbumDates(albumId: AlbumId, startDate: Date, endDate: Date): Promise<void>;
    fetchAlbums(): Promise<Album[]>;
    fetchMedias(albumId: AlbumId): Promise<Media[]>;
}

export async function updateAlbumDatesThunk(
    dispatch: (action: LoadingStarted | AlbumDatesUpdated) => void,
    updateAlbumDatesPort: UpdateAlbumDatesPort,
    editDatesDialogState: EditDatesDialogState
): Promise<void> {
    dispatch(loadingStarted());

    await updateAlbumDatesPort.updateAlbumDates(
        editDatesDialogState.albumId,
        editDatesDialogState.startDate,
        editDatesDialogState.endDate
    );

    const [albums, medias] = await Promise.all([
        updateAlbumDatesPort.fetchAlbums(),
        updateAlbumDatesPort.fetchMedias(editDatesDialogState.albumId)
    ]);

    dispatch(albumDatesUpdated({
        albums,
        medias,
        updatedAlbumId: editDatesDialogState.albumId
    }));
}

export const updateAlbumDatesDeclaration: ThunkDeclaration<
    CatalogViewerState,
    { editDatesDialogState?: EditDatesDialogState },
    () => Promise<void>,
    CatalogFactoryArgs
> = {
    selector: ({editDatesDialog}: CatalogViewerState) => ({
        editDatesDialogState: editDatesDialog,
    }),

    factory: ({dispatch, app, partialState: {editDatesDialogState}}) => {
        if (!editDatesDialogState) {
            return async () => {};
        }

        const updateAlbumDatesPort: UpdateAlbumDatesPort = new CatalogAPIAdapter(app.axiosInstance, app);
        return updateAlbumDatesThunk.bind(null, dispatch, updateAlbumDatesPort, editDatesDialogState);
    },
};
