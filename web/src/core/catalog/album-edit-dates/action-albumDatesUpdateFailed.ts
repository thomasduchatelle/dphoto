import {CatalogViewerState} from "../language";
import {createAction} from "src/libs/daction";

export const ALBUM_DATES_ORPHANED_MEDIAS_ERROR_CODE = "OrphanedMediasError";

interface AlbumDatesUpdateFailedPayload {
    error?: string; // This will be the error code or message from the API
}

export const albumDatesUpdateFailed = createAction<CatalogViewerState, AlbumDatesUpdateFailedPayload>(
    "AlbumDatesUpdateFailed",
    (current: CatalogViewerState, {error}: AlbumDatesUpdateFailedPayload) => {
        if (!current.editDatesDialog) {
            return current;
        }

        let errorMessage: string;
        switch (error) {
            case ALBUM_DATES_ORPHANED_MEDIAS_ERROR_CODE:
                errorMessage = "The dates cannot be updated if it makes some pictures orphaned. Create another album first.";
                break;
            default:
                errorMessage = error || "Album dates couldn't be saved. Refresh your page and retry, or let the developer known.";
                break;
        }

        return {
            ...current,
            editDatesDialog: {
                ...current.editDatesDialog,
                isLoading: false,
                error: errorMessage, // Store the user-friendly message
            },
        };
    }
);

export type AlbumDatesUpdateFailed = ReturnType<typeof albumDatesUpdateFailed>;
