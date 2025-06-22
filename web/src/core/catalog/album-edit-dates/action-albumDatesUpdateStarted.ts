import {CatalogViewerState} from "../language";
import {createAction} from "../../../light-state-lib/action-factory";

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
            },
        };
    }
);

export type AlbumDatesUpdateStarted = ReturnType<typeof albumDatesUpdateStarted>;
