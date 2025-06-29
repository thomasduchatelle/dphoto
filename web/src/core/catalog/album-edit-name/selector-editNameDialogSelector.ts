import {CatalogViewerState, EditNameDialog, isEditNameDialog} from "../language";
import {baseEditNameSelector, BaseEditNameSelection} from "../base-name-edit/selector-baseEditNameSelector";

export interface EditNameDialogSelection extends BaseEditNameSelection {
    isOpen: boolean;
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
    const baseSelection = baseEditNameSelector(state);
    
    if (!baseSelection) {
        return closedEditNameSelection;
    }

    const dialog: EditNameDialog = state.dialog as EditNameDialog;
    const isSaveEnabled = !baseSelection.nameError && !baseSelection.folderNameError && !dialog.isLoading;

    return {
        ...baseSelection,
        isOpen: true,
        isLoading: dialog.isLoading,
        isSaveEnabled,
    };
}
