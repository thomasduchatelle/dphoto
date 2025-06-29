import {CatalogViewerState, isEditNameDialog} from "../language";
import {BaseEditNameSelection, baseEditNameSelector} from "../base-name-edit/selector-baseEditNameSelector";

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
    if (!isEditNameDialog(state.dialog)) {
        return closedEditNameSelection;
    }

    const {isSavable, ...baseSelection} = baseEditNameSelector(state, state.dialog);

    if (!baseSelection) {
        return closedEditNameSelection;
    }

    const isSaveEnabled = isSavable && !state.dialog.isLoading;

    return {
        ...baseSelection,
        isOpen: true,
        isLoading: state.dialog.isLoading,
        isSaveEnabled,
    };
}
