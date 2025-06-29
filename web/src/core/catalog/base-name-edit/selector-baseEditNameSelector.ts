import {albumIdEquals, CatalogViewerState, EditNameDialog} from "../language";

export interface BaseEditNameSelection {
    albumName: string;
    originalName: string;
    customFolderName: string;
    isCustomFolderNameEnabled: boolean;
    technicalError?: string;
    nameError?: string;
    folderNameError?: string;
}

export function baseEditNameSelector(state: CatalogViewerState, dialog: EditNameDialog): BaseEditNameSelection & { isSavable: boolean } {
    const {technicalError, nameError, folderNameError} = dialog.nameError

    const originalName = state.allAlbums.find(album => albumIdEquals(album.albumId, dialog.albumId))?.name || dialog.albumName;

    return {
        albumName: dialog.albumName,
        originalName,
        customFolderName: dialog.customFolderName,
        isCustomFolderNameEnabled: dialog.isCustomFolderNameEnabled,
        technicalError,
        nameError,
        folderNameError,
        isSavable: !!nameError && !!folderNameError,
    };
}
