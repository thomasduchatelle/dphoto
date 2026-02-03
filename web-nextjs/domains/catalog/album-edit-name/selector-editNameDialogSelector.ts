import {albumIdEquals, CatalogViewerState, EditNameDialog, isEditNameDialog} from "../language";
import {BaseEditNameSelection, baseEditNameSelector} from "../base-name-edit/selector-baseEditNameSelector";

export interface EditNameDialogSelection extends BaseEditNameSelection {
    isOpen: boolean;
    originalName: string;
    technicalError?: string;
    isLoading: boolean;
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

    const {isSavable, ...baseSelection} = baseEditNameSelector(state, dialog);

    if (!baseSelection) {
        return closedEditNameSelection;
    }

    const isSaveEnabled = isSavable && !dialog.isLoading;

    const originalName = state.allAlbums.find(album => albumIdEquals(album.albumId, dialog.albumId))?.name || "";

    return {
        ...baseSelection,
        isOpen: true,
        isLoading: state.dialog.isLoading,
        technicalError: state.dialog.technicalError,
        isSaveEnabled,
        originalName,
    };
}
