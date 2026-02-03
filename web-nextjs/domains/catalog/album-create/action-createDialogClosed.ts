import {createAction} from "@/libs/daction";
import {CatalogViewerState} from "../language";

export const createDialogClosed = createAction<CatalogViewerState>(
    "CreateDialogClosed",
    (current: CatalogViewerState) => {
        return {
            ...current,
            dialog: undefined,
        };
    }
);

export type CreateDialogClosed = ReturnType<typeof createDialogClosed>;
