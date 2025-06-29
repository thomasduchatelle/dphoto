import {createAction} from "src/libs/daction";
import {CatalogViewerState, isNameEditBase} from "../language";

export const customFolderNameToggled = createAction<CatalogViewerState, boolean>(
    "CustomFolderNameToggled",
    (current: CatalogViewerState, isFolderNameEnabled: boolean) => {
        if (!isNameEditBase(current.dialog)) {
            return current;
        }

        return {
            ...current,
            dialog: {
                ...current.dialog,
                isCustomFolderNameEnabled: isFolderNameEnabled,
                customFolderName: isFolderNameEnabled ? (current.dialog.originalFolderName ?? "") : "",
                nameError: {
                    ...current.dialog.nameError,
                    folderNameError: undefined,
                },
            },
        };
    }
);
