import {AlbumId, Album, Media, MediaWithinADay} from "../language";
import {albumDatesUpdateStarted, AlbumDatesUpdateStarted} from "./action-albumDatesUpdateStarted";
import {albumDatesUpdated, AlbumDatesUpdated} from "./action-albumDatesUpdated";
import {ThunkDeclaration} from "../../thunk-engine";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {CatalogAPIAdapter} from "../adapters/api";
import {MediaPerDayLoader} from "../navigation/MediaPerDayLoader";

export interface UpdateAlbumDatesPort {
    updateAlbumDates(albumId: AlbumId, startDate: Date, endDate: Date): Promise<void>;
    fetchAlbums(): Promise<Album[]>;
    fetchMedias(albumId: AlbumId): Promise<{medias: MediaWithinADay[]}>;
}

export async function updateAlbumDatesThunk(
    dispatch: (action: AlbumDatesUpdateStarted | AlbumDatesUpdated) => void,
    updateAlbumDatesPort: UpdateAlbumDatesPort,
    albumId: AlbumId,
    startDate: Date,
    endDate: Date
): Promise<void> {
    dispatch(albumDatesUpdateStarted());

    const apiStartDate = new Date(startDate);
    apiStartDate.setHours(0, 0, 0, 0);

    const apiEndDate = new Date(endDate);
    apiEndDate.setDate(apiEndDate.getDate() + 1);
    apiEndDate.setHours(0, 0, 0, 0);

    await updateAlbumDatesPort.updateAlbumDates(albumId, apiStartDate, apiEndDate);

    const [albums, mediasResp] = await Promise.all([
        updateAlbumDatesPort.fetchAlbums(),
        updateAlbumDatesPort.fetchMedias(albumId)
    ]);

    dispatch(albumDatesUpdated({albums, medias: mediasResp.medias}));
}

export const updateAlbumDatesDeclaration: ThunkDeclaration<
    any,
    {albumId: AlbumId, startDate: Date, endDate: Date},
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
        const updateAlbumDatesPort: UpdateAlbumDatesPort = {
            ...restAdapter,
            fetchMedias: (albumId: AlbumId) => new MediaPerDayLoader(restAdapter).loadMedias(albumId)
        };
        return (newStartDate: Date, newEndDate: Date) => 
            updateAlbumDatesThunk(dispatch, updateAlbumDatesPort, albumId, newStartDate, newEndDate);
    },
};
