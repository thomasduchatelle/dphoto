import {CatalogViewerState, isEditDatesDialog} from "../language";
import {createAction} from "src/libs/daction";

export const ALBUM_DATES_ORPHANED_MEDIAS_ERROR_CODE = "OrphanedMediasError";

interface AlbumDatesUpdateFailedPayload {
    error?: string;
}

export const albumDatesUpdateFailed = createAction<CatalogViewerState, AlbumDatesUpdateFailedPayload>(
    "AlbumDatesUpdateFailed",
    (current: CatalogViewerState, {error}: AlbumDatesUpdateFailedPayload) => {
        if (!isEditDatesDialog(current.dialog)) {
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
            dialog: {
                ...current.dialog,
                isLoading: false,
                error: errorMessage,
            },
        };
    }
);

export type AlbumDatesUpdateFailed = ReturnType<typeof albumDatesUpdateFailed>;
