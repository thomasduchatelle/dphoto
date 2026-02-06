import {CatalogViewerState, NameEditBase} from "../language";

export interface BaseEditNameSelection {
    albumName: string;
    customFolderName: string;
    isCustomFolderNameEnabled: boolean;
    nameError?: string;
    folderNameError?: string;
}

export interface BaseEditNameSelectionWithSavable extends BaseEditNameSelection {
    isSavable: boolean
}

export function baseEditNameSelector(state: CatalogViewerState, dialog: NameEditBase): BaseEditNameSelectionWithSavable {
    const {nameError, folderNameError} = dialog.nameError

    return {
        albumName: dialog.albumName,
        customFolderName: dialog.customFolderName,
        isCustomFolderNameEnabled: dialog.isCustomFolderNameEnabled,
        nameError,
        folderNameError,
        isSavable: !!dialog.albumName && (!dialog.isCustomFolderNameEnabled || !!dialog.customFolderName) && !nameError && !folderNameError,
    };
}
