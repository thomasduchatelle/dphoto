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

export interface UpdateAlbumDatesThunkArgs {
    albumId: AlbumId;
    startDate: Date;
    endDate: Date;
}

export async function updateAlbumDatesThunk(
    dispatch: (action: AlbumDatesUpdateStarted | AlbumDatesUpdated) => void,
    updateAlbumDatesPort: UpdateAlbumDatesPort,
    mediaPerDayLoader: MediaPerDayLoader,
    args: UpdateAlbumDatesThunkArgs
): Promise<void> {
    dispatch(albumDatesUpdateStarted());

    const apiStartDate = new Date(Date.UTC(args.startDate.getUTCFullYear(), args.startDate.getUTCMonth(), args.startDate.getUTCDate()));

    const apiEndDate = new Date(Date.UTC(args.endDate.getUTCFullYear(), args.endDate.getUTCMonth(), args.endDate.getUTCDate() + 1));

    await updateAlbumDatesPort.updateAlbumDates(args.albumId, apiStartDate, apiEndDate);

    const [albums, mediasResp] = await Promise.all([
        updateAlbumDatesPort.fetchAlbums(),
        mediaPerDayLoader.loadMedias(args.albumId)
    ]);

    dispatch(albumDatesUpdated({albums, medias: mediasResp.medias}));
}

export const updateAlbumDatesDeclaration: ThunkDeclaration<
    any,
    { albumId: AlbumId, startDate: Date, endDate: Date },
    (args: UpdateAlbumDatesThunkArgs) => Promise<void>,
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
        return (args: UpdateAlbumDatesThunkArgs) =>
            updateAlbumDatesThunk(dispatch, updateAlbumDatesPort, mediaPerDayLoader, args);
    },
};
