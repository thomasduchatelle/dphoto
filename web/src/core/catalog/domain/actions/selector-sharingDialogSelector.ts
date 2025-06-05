import {CatalogViewerState, ShareError, Sharing} from "../catalog-state";

export interface SharingDialogFrag {
    open: boolean;
    sharedWith: Sharing[];
    error?: ShareError;
}

export function sharingDialogSelector({shareModal}: CatalogViewerState): SharingDialogFrag {
    if (!shareModal) {
        return {
            open: false,
            sharedWith: [],
        };
    }
    return {
        open: true,
        sharedWith: shareModal.sharedWith,
        error: shareModal.error,
    };
}
