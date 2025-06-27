import {createAction} from "src/libs/daction";
import {CatalogViewerState, editNameDialogNoError, isEditNameDialog} from "../language";

export const albumNameChanged = createAction<CatalogViewerState, string>(
    "AlbumNameChanged",
    (current: CatalogViewerState, albumName: string) => {
        if (!isEditNameDialog(current.dialog)) {
            return current;
        }


        return {
            ...current,
            dialog: {
                ...current.dialog,
                albumName,
                error: !!albumName ? editNameDialogNoError : {nameError: "Album name is mandatory"},
            },
        };
    }
);

export type AlbumNameChanged = ReturnType<typeof albumNameChanged>;
