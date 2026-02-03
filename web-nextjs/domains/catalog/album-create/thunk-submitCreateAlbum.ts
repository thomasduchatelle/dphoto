import {ThunkDeclaration} from "@/libs/dthunks";
import {Album, AlbumId, CatalogViewerAction, CatalogViewerState, getErrorMessage, isCatalogError, isCreateDialog} from "../language";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {CatalogFactory} from "../catalog-factories";
import {DPhotoApplication} from "../../application";
import {createAlbumStarted} from "./action-createAlbumStarted";
import {createAlbumFailed} from "./action-createAlbumFailed";
import {albumsLoaded} from "../navigation";
import {convertToModelEndDate, convertToModelStartDate, validateDateRange} from "../date-range/date-helper";

interface CreateDialogData {
    albumName: string;
    startDate: Date | null;
    endDate: Date | null;
    customFolderName: string;
    startAtDayStart: boolean;
    endAtDayEnd: boolean;
    isCustomFolderNameEnabled: boolean;
}

export interface CreateAlbumRequest {
    name: string
    start: Date
    end: Date
    forcedFolderName: string
}

export interface CreateAlbumPort {
    createAlbum(request: CreateAlbumRequest): Promise<AlbumId>

    fetchAlbums(): Promise<Album[]>
}

export async function submitCreateAlbumThunk(
    dispatch: (action: CatalogViewerAction) => void,
    createAlbumPort: CreateAlbumPort,
    dialogData: CreateDialogData
): Promise<void> {
    const {albumName, startDate, endDate, customFolderName, startAtDayStart, endAtDayEnd, isCustomFolderNameEnabled} = dialogData;

    const validation = validateDateRange({
        startDate,
        endDate,
        startAtDayStart,
        endAtDayEnd,
    });

    if (!validation.areDatesValid || !validation.isDateRangeValid) {
        dispatch(createAlbumFailed("AlbumStartAndEndDateMandatoryErr"));
        return;
    }

    dispatch(createAlbumStarted());

    try {
        const actualStart = convertToModelStartDate(startDate!, startAtDayStart);
        const actualEnd = convertToModelEndDate(endDate!, endAtDayEnd);

        const albumId = await createAlbumPort.createAlbum({
            name: albumName,
            start: actualStart,
            end: actualEnd,
            forcedFolderName: isCustomFolderNameEnabled ? customFolderName : "",
        });

        const albums = await createAlbumPort.fetchAlbums();
        dispatch(albumsLoaded({albums, redirectTo: albumId}));
    } catch (error) {
        const errorMessage = isCatalogError(error) ? error.code : getErrorMessage(error) ?? "Unknown error";
        dispatch(createAlbumFailed(errorMessage));
        throw error;
    }
}

export const submitCreateAlbumDeclaration: ThunkDeclaration<
    CatalogViewerState,
    CreateDialogData,
    () => Promise<void>,
    CatalogFactoryArgs
> = {
    selector: (state: CatalogViewerState) => {
        const dialog = state.dialog;
        if (!isCreateDialog(dialog)) {
            return {
                albumName: "",
                startDate: null,
                endDate: null,
                customFolderName: "",
                startAtDayStart: true,
                endAtDayEnd: true,
                isCustomFolderNameEnabled: false,
            };
        }

        return {
            albumName: dialog.albumName,
            startDate: dialog.startDate,
            endDate: dialog.endDate,
            customFolderName: dialog.customFolderName,
            startAtDayStart: dialog.startAtDayStart,
            endAtDayEnd: dialog.endAtDayEnd,
            isCustomFolderNameEnabled: dialog.isCustomFolderNameEnabled,
        };
    },
    factory: ({dispatch, app, partialState}) => {
        const restAdapter = new CatalogFactory(app as DPhotoApplication).restAdapter();
        return submitCreateAlbumThunk.bind(null, dispatch, restAdapter, partialState);
    },
};
