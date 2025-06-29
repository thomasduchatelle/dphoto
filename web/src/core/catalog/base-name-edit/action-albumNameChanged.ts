import {createAction} from "src/libs/daction";
import {CatalogViewerState, isNameEditBase} from "../language";

export const albumNameChanged = createAction<CatalogViewerState, string>(
    "AlbumNameChanged",
    (current: CatalogViewerState, albumName: string) => {
        if (!isNameEditBase(current.dialog)) {
            return current;
        }


        return {
            ...current,
            dialog: {
                ...current.dialog,
                albumName,
                nameError: {
                    ...current.dialog.nameError,
                    nameError: !!albumName ? undefined : "Album name is mandatory",
                },
            },
        };
    }
);

export type AlbumNameChanged = ReturnType<typeof albumNameChanged>;
