import {Album, AlbumId, CatalogViewerState} from "../language";
import {albumDatesUpdateStarted} from "./action-albumDatesUpdateStarted";
import {albumDatesUpdated} from "./action-albumDatesUpdated";
import {ThunkDeclaration} from "../../thunk-engine";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {CatalogAPIAdapter} from "../adapters/api";
import {MediaPerDayLoader} from "../navigation";
import {Action} from "src/light-state-lib";

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
    dispatch: (action: Action<CatalogViewerState, any>) => void,
    updateAlbumDatesPort: UpdateAlbumDatesPort,
    mediaPerDayLoader: MediaPerDayLoader,
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

    const [albums, mediasResp] = await Promise.all([
        updateAlbumDatesPort.fetchAlbums(),
        mediaPerDayLoader.loadMedias(albumId)
    ]);

    dispatch(albumDatesUpdated({albums, medias: mediasResp.medias}));
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
        const restAdapter = new CatalogAPIAdapter(app.axiosInstance, app);
        const mediaPerDayLoader = new MediaPerDayLoader(restAdapter);
        const updateAlbumDatesPort: UpdateAlbumDatesPort = restAdapter;
        return () =>
            updateAlbumDatesThunk(dispatch, updateAlbumDatesPort, mediaPerDayLoader, partialState);
    },
};
