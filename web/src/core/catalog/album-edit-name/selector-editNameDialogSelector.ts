import {albumIdEquals, CatalogViewerState, EditNameDialog, isEditNameDialog} from "../language";

export interface EditNameDialogSelection {
    isOpen: boolean;
    albumName: string;
    originalName: string;
    customFolderName: string;
    isCustomFolderNameEnabled: boolean;
    isLoading: boolean;
    technicalError?: string;
    nameError?: string;
    folderNameError?: string;
    isSaveEnabled: boolean;
}

const closedEditNameSelection: EditNameDialogSelection = {
    isOpen: false,
    albumName: "",
    originalName: "",
    customFolderName: "",
    isCustomFolderNameEnabled: false,
    isLoading: false,
    isSaveEnabled: false,
};

export function editNameDialogSelector(state: CatalogViewerState): EditNameDialogSelection {
    if (!isEditNameDialog(state.dialog)) {
        return closedEditNameSelection;
    }

    const dialog: EditNameDialog = state.dialog;
    const {technicalError, nameError, folderNameError} = dialog.error

    const originalName = state.allAlbums.find(album => albumIdEquals(album.albumId, dialog.albumId))?.name || dialog.albumName;

    const isSaveEnabled = !nameError && !folderNameError && !dialog.isLoading;

    return {
        isOpen: true,
        albumName: dialog.albumName,
        originalName,
        customFolderName: dialog.customFolderName,
        isCustomFolderNameEnabled: dialog.isCustomFolderNameEnabled,
        isLoading: dialog.isLoading,
        technicalError,
        nameError,
        folderNameError,
        isSaveEnabled,
    };
}
