import {ThunkDeclaration} from "src/libs/dthunks";
import {CatalogViewerAction, CatalogViewerState, getErrorMessage, isCatalogError, isCreateDialog} from "../language";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {CreateAlbumPort} from "./album-createAlbum";
import {CatalogFactory} from "../catalog-factories";
import {DPhotoApplication} from "../../application";
import {createAlbumStarted} from "./action-createAlbumStarted";
import {createAlbumFailed} from "./action-createAlbumFailed";
import {albumsLoaded} from "../navigation";
import {convertToModelEndDate, convertToModelStartDate, validateDateRange} from "../date-range/date-helper";

interface CreateDialogData {
    name: string;
    startDate: Date | null;
    endDate: Date | null;
    forceFolderName: string;
    startAtDayStart: boolean;
    endAtDayEnd: boolean;
    withCustomFolderName: boolean;
}

export async function submitCreateAlbumThunk(
    dispatch: (action: CatalogViewerAction) => void,
    createAlbumPort: CreateAlbumPort,
    dialogData: CreateDialogData
): Promise<void> {
    const {name, startDate, endDate, forceFolderName, startAtDayStart, endAtDayEnd, withCustomFolderName} = dialogData;
    
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
            name,
            start: actualStart,
            end: actualEnd,
            forcedFolderName: withCustomFolderName ? forceFolderName : "",
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
                name: "",
                startDate: null,
                endDate: null,
                forceFolderName: "",
                startAtDayStart: true,
                endAtDayEnd: true,
                withCustomFolderName: false,
            };
        }

        return {
            name: dialog.name,
            startDate: dialog.startDate,
            endDate: dialog.endDate,
            forceFolderName: dialog.forceFolderName,
            startAtDayStart: dialog.startAtDayStart,
            endAtDayEnd: dialog.endAtDayEnd,
            withCustomFolderName: dialog.withCustomFolderName,
        };
    },
    factory: ({dispatch, app, partialState}) => {
        const restAdapter = new CatalogFactory(app as DPhotoApplication).restAdapter();
        return submitCreateAlbumThunk.bind(null, dispatch, restAdapter, partialState);
    },
};
