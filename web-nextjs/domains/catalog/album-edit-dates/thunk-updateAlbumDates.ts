import {Album, AlbumId, CatalogViewerState, getErrorMessage, isCatalogError, isEditDatesDialog, Media} from "../language";
import {albumDatesUpdateStarted} from "./action-albumDatesUpdateStarted";
import {albumsAndMediasLoaded} from "@/domains/catalog";
import {albumDatesUpdateFailed} from "./action-albumDatesUpdateFailed";
import {Action} from "@/libs/daction";
import {ThunkDeclaration} from "@/libs/dthunks";
import {convertToModelEndDate, convertToModelStartDate} from "../date-range/date-helper";
import {CatalogDispatch} from "@/domains/catalog/common/catalog-dispatch";

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
        const apiStartDate = convertToModelStartDate(startDate, startAtDayStart);
        const apiEndDate = convertToModelEndDate(endDate, endAtDayEnd);

        await updateAlbumDatesPort.updateAlbumDates(albumId, apiStartDate, apiEndDate);

        const [albums, medias] = await Promise.all([
            updateAlbumDatesPort.fetchAlbums(),
            updateAlbumDatesPort.fetchMedias(albumId)
        ]);

        dispatch(albumsAndMediasLoaded({albums, medias, mediasFromAlbumId: albumId}));

    } catch (error) {
        if (isCatalogError(error) && error.code === editDatesOrphanedMediasErrorCode) {
            dispatch(albumDatesUpdateFailed({error: error.message}));
        } else {
            dispatch(albumDatesUpdateFailed({error: getErrorMessage(error)}));
        }
    }
}


export const updateAlbumDatesDeclaration: ThunkDeclaration<
    CatalogViewerState,
    UpdateAlbumDatesThunkArgs | undefined,
    () => Promise<void>,
    CatalogDispatch & { adapter: UpdateAlbumDatesPort }
> = {
    selector: (state: CatalogViewerState) => {
        const dialog = state.dialog;
        if (!isEditDatesDialog(dialog) || !dialog.startDate || !dialog.endDate) {
            return undefined;
        }
        return {
            albumId: dialog.albumId,
            startDate: dialog.startDate,
            endDate: dialog.endDate,
            startAtDayStart: dialog.startAtDayStart,
            endAtDayEnd: dialog.endAtDayEnd,
        };
    },

    factory: ({dispatch, adapter, partialState}) => {
        return () =>
            updateAlbumDatesThunk(dispatch, adapter, partialState);
    },
};
