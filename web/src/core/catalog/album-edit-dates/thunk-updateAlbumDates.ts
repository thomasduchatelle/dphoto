import {Album, AlbumId} from "../language";
import {albumDatesUpdateStarted, AlbumDatesUpdateStarted} from "./action-albumDatesUpdateStarted";
import {albumDatesUpdated, AlbumDatesUpdated} from "./action-albumDatesUpdated";
import {ThunkDeclaration} from "../../thunk-engine";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {CatalogAPIAdapter} from "../adapters/api";
import {MediaPerDayLoader} from "../navigation";

export interface UpdateAlbumDatesPort {
    updateAlbumDates(albumId: AlbumId, startDate: Date, endDate: Date): Promise<void>;

    fetchAlbums(): Promise<Album[]>;
}

export async function updateAlbumDatesThunk(
    dispatch: (action: AlbumDatesUpdateStarted | AlbumDatesUpdated) => void,
    updateAlbumDatesPort: UpdateAlbumDatesPort,
    mediaPerDayLoader: MediaPerDayLoader,
    albumId: AlbumId,
    startDate: Date,
    endDate: Date
): Promise<void> {
    dispatch(albumDatesUpdateStarted());

    const apiStartDate = new Date(Date.UTC(startDate.getUTCFullYear(), startDate.getUTCMonth(), startDate.getUTCDate()));

    const apiEndDate = new Date(Date.UTC(endDate.getUTCFullYear(), endDate.getUTCMonth(), endDate.getUTCDate() + 1));

    await updateAlbumDatesPort.updateAlbumDates(albumId, apiStartDate, apiEndDate);

    const [albums, mediasResp] = await Promise.all([
        updateAlbumDatesPort.fetchAlbums(),
        mediaPerDayLoader.loadMedias(albumId)
    ]);

    dispatch(albumDatesUpdated({albums, medias: mediasResp.medias}));
}

export const updateAlbumDatesDeclaration: ThunkDeclaration<
    any,
    { albumId: AlbumId, startDate: Date, endDate: Date },
    (startDate: Date, endDate: Date) => Promise<void>,
    CatalogFactoryArgs
> = {
    selector: (state: any) => ({
        albumId: state.editDatesDialog?.albumId,
        startDate: state.editDatesDialog?.startDate,
        endDate: state.editDatesDialog?.endDate,
    }),

    factory: ({dispatch, app, partialState: {albumId, startDate, endDate}}) => {
        const restAdapter = new CatalogAPIAdapter(app.axiosInstance, app);
        const mediaPerDayLoader = new MediaPerDayLoader(restAdapter); // Instantiate MediaPerDayLoader
        const updateAlbumDatesPort: UpdateAlbumDatesPort = restAdapter; // Use restAdapter directly for the port
        return (newStartDate: Date, newEndDate: Date) =>
            updateAlbumDatesThunk(dispatch, updateAlbumDatesPort, mediaPerDayLoader, albumId, newStartDate, newEndDate);
    },
};
