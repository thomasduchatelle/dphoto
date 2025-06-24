import {Album, AlbumId, CatalogViewerState, getErrorMessage, isCatalogError, Media} from "../language";
import {albumDatesUpdateStarted} from "./action-albumDatesUpdateStarted";
import {albumDatesUpdated} from "./action-albumDatesUpdated";
import {albumDatesUpdateFailed} from "./action-albumDatesUpdateFailed";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {CatalogAPIAdapter} from "../adapters/api";
import {Action} from "src/libs/daction";
import {groupByDay} from "../navigation/group-by-day";
import {ThunkDeclaration} from "src/libs/dthunks";
import {isRoundTime} from "../common/date-helper";

/** When deletion or date edit is not possible because it would orphan medias */
export const editDatesOrphanedMediasErrorCode = "OrphanedMediasErr";

export interface UpdateAlbumDatesPort {
    updateAlbumDates(albumId: AlbumId, startDate: Date, endDate: Date): Promise<void>;

    fetchAlbums(): Promise<Album[]>;

    fetchMedias(albumId: AlbumId): Promise<Media[]>;
}

export interface UpdateAlbumDatesThunkArgs {
    albumId: AlbumId;
    startDate: Date;
    endDate: Date;
    startAtDayStart: boolean;
    endAtDayEnd: boolean;
}

export async function updateAlbumDatesThunk(
    dispatch: (action: Action<CatalogViewerState, any>) => void,
    updateAlbumDatesPort: UpdateAlbumDatesPort,
    dialog?: UpdateAlbumDatesThunkArgs
): Promise<void> {
    if (!dialog) {
        return
    }

    const {albumId, startDate, endDate, startAtDayStart, endAtDayEnd} = dialog;

    dispatch(albumDatesUpdateStarted());

    try {
        const apiStartDate = convertToApiStartDate(startDate, startAtDayStart);
        const apiEndDate = convertToApiEndDate(endDate, endAtDayEnd);

        await updateAlbumDatesPort.updateAlbumDates(albumId, apiStartDate, apiEndDate);

        const [albums, medias] = await Promise.all([
            updateAlbumDatesPort.fetchAlbums(),
            updateAlbumDatesPort.fetchMedias(albumId)
        ]);

        dispatch(albumDatesUpdated({albums, medias: groupByDay(medias)}));

    } catch (error) {
        if (isCatalogError(error) && error.code === editDatesOrphanedMediasErrorCode) {
            dispatch(albumDatesUpdateFailed({error: error.message}));
        } else {
            dispatch(albumDatesUpdateFailed({error: getErrorMessage(error)}));
        }
    }
}

function convertToApiStartDate(original: Date, atDayStart: boolean): Date {
    const date = new Date(original);
    if (atDayStart) {
        return new Date(Date.UTC(date.getUTCFullYear(), date.getUTCMonth(), date.getUTCDate()));
    }
    
    return new Date(Date.UTC(date.getUTCFullYear(), date.getUTCMonth(), date.getUTCDate(), date.getHours(), date.getMinutes()));
}

function convertToApiEndDate(original: Date, atDayEnd: boolean): Date {
    const date = new Date(original);
    if (atDayEnd) {
        return new Date(Date.UTC(date.getUTCFullYear(), date.getUTCMonth(), date.getUTCDate() + 1));
    }

    if (isRoundTime(date)) {
        return new Date(Date.UTC(date.getUTCFullYear(), date.getUTCMonth(), date.getUTCDate(), date.getHours(), date.getMinutes()));
    }

    return new Date(Date.UTC(date.getUTCFullYear(), date.getUTCMonth(), date.getUTCDate(), date.getHours(), date.getMinutes() + 1));
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
            startAtDayStart: state.editDatesDialog.startAtDayStart,
            endAtDayEnd: state.editDatesDialog.endAtDayEnd,
        } : undefined
    ),

    factory: ({dispatch, app, partialState}) => {
        const updateAlbumDatesPort: UpdateAlbumDatesPort = new CatalogAPIAdapter(app.axiosInstance, app);
        return () =>
            updateAlbumDatesThunk(dispatch, updateAlbumDatesPort, partialState);
    },
};
