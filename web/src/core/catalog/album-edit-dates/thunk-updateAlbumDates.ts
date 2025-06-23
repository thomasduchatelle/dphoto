import {Album, AlbumId, CatalogViewerState, Media} from "../language";
import {albumDatesUpdateStarted} from "./action-albumDatesUpdateStarted";
import {albumDatesUpdated} from "./action-albumDatesUpdated";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {CatalogAPIAdapter} from "../adapters/api";
import {Action} from "src/libs/daction";
import {groupByDay} from "../navigation/group-by-day";
import {ThunkDeclaration} from "src/libs/dthunks";

export interface UpdateAlbumDatesPort {
    updateAlbumDates(albumId: AlbumId, startDate: Date, endDate: Date): Promise<void>;

    fetchAlbums(): Promise<Album[]>;

    fetchMedias(albumId: AlbumId): Promise<Media[]>;
}

export interface UpdateAlbumDatesThunkArgs {
    albumId: AlbumId;
    startDate: Date;
    endDate: Date;
}

export async function updateAlbumDatesThunk(
    dispatch: (action: Action<CatalogViewerState, any>) => void,
    updateAlbumDatesPort: UpdateAlbumDatesPort,
    dialog?: UpdateAlbumDatesThunkArgs
): Promise<void> {
    if (!dialog) {
        return
    }

    let {albumId, startDate, endDate} = dialog;
    startDate = new Date(startDate);
    endDate = new Date(endDate);

    dispatch(albumDatesUpdateStarted());

    const apiStartDate = new Date(Date.UTC(startDate.getUTCFullYear(), startDate.getUTCMonth(), startDate.getUTCDate()));
    const apiEndDate = new Date(Date.UTC(endDate.getUTCFullYear(), endDate.getUTCMonth(), endDate.getUTCDate() + 1));

    await updateAlbumDatesPort.updateAlbumDates(albumId, apiStartDate, apiEndDate);

    const [albums, medias] = await Promise.all([
        updateAlbumDatesPort.fetchAlbums(),
        updateAlbumDatesPort.fetchMedias(albumId)
    ]);

    dispatch(albumDatesUpdated({albums, medias: groupByDay(medias)}));
}

export const updateAlbumDatesDeclaration: ThunkDeclaration<
    CatalogViewerState,
    UpdateAlbumDatesThunkArgs | undefined,
    () => Promise<void>,
    CatalogFactoryArgs
> = {
    selector: (state: CatalogViewerState) => (state.editDatesDialog ? {
            albumId: state.editDatesDialog.albumId,
            startDate: state.editDatesDialog.startDate,
            endDate: state.editDatesDialog.endDate,
        } : undefined
    ),

    factory: ({dispatch, app, partialState}) => {
        const updateAlbumDatesPort: UpdateAlbumDatesPort = new CatalogAPIAdapter(app.axiosInstance, app);
        return () =>
            updateAlbumDatesThunk(dispatch, updateAlbumDatesPort, partialState);
    },
};
