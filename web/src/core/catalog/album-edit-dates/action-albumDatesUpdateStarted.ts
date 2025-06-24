import {CatalogViewerState, isEditDatesDialog} from "../language";
import {createAction} from "src/libs/daction";

export const albumDatesUpdateStarted = createAction<CatalogViewerState>(
    "AlbumDatesUpdateStarted",
    (current: CatalogViewerState) => {
        const dialog = current.dialog;
        if (!isEditDatesDialog(dialog)) {
            return current;
        }

        return {
            ...current,
            dialog: {
                ...dialog,
                isLoading: true,
                error: undefined,
            },
        };
    }
);

export type AlbumDatesUpdateStarted = ReturnType<typeof albumDatesUpdateStarted>;
