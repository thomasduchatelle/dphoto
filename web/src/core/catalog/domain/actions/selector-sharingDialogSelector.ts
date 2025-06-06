import {CatalogViewerState, ShareError, Sharing, UserDetails} from "../catalog-state";

export interface SharingDialogFrag {
    open: boolean;
    sharedWith: Sharing[];
    suggestions: UserDetails[];
    error?: ShareError;
}

export function sharingDialogSelector({shareModal}: CatalogViewerState): SharingDialogFrag {
    if (!shareModal) {
        return {
            open: false,
            sharedWith: [],
            suggestions: [],
        };
    }
    return {
        open: true,
        sharedWith: shareModal.sharedWith,
        suggestions: shareModal.suggestions,
        error: shareModal.error,
    };
}
