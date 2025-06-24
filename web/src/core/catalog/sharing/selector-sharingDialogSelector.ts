import {CatalogViewerState, ShareError, Sharing, UserDetails, isShareDialog} from "../language";

export interface SharingDialogFrag {
    open: boolean;
    sharedWith: Sharing[];
    suggestions: UserDetails[];
    error?: ShareError;
}

export function sharingDialogSelector(state: CatalogViewerState): SharingDialogFrag {
    if (!isShareDialog(state.dialog)) {
        return {
            open: false,
            sharedWith: [],
            suggestions: [],
        };
    }
    return {
        open: true,
        sharedWith: state.dialog.sharedWith,
        suggestions: state.dialog.suggestions,
        error: state.dialog.error,
    };
}
