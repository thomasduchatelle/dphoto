import {CatalogViewerState} from "../language";
import {createAction} from "src/libs/daction";

export const albumDatesUpdateStarted = createAction<CatalogViewerState>(
    "AlbumDatesUpdateStarted",
    (current: CatalogViewerState) => {
        if (!current.editDatesDialog) {
            return current;
        }

        return {
            ...current,
            editDatesDialog: {
                ...current.editDatesDialog,
                isLoading: true,
                error: undefined,
            },
        };
    }
);

export type AlbumDatesUpdateStarted = ReturnType<typeof albumDatesUpdateStarted>;
