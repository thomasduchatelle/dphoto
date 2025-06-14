import {CatalogViewerState} from "../language";
import {albumIdEquals} from "../language/utils-albumIdEquals";

export interface EditDatesDialogProperties {
    isOpen: boolean;
    albumName: string;
}

export const closedEditDatesDialogProperties: EditDatesDialogProperties = {
    isOpen: false,
    albumName: ""
};

export function selectEditDatesDialog(state: CatalogViewerState): EditDatesDialogProperties {
    if (!state.editDatesDialog) {
        return closedEditDatesDialogProperties;
    }

    const album = state.allAlbums.find(album => 
        albumIdEquals(album.albumId, state.editDatesDialog?.albumId)
    );

    if (!album) {
        return closedEditDatesDialogProperties;
    }

    return {
        isOpen: true,
        albumName: album.name
    };
}
