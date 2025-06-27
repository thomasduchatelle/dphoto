import {createAction} from "src/libs/daction";
import {albumIdEquals, CatalogViewerState} from "../language";
import {getDisplayedAlbumId} from "../language/selector-displayedAlbum";

export const editNameDialogOpened = createAction<CatalogViewerState>(
    "EditNameDialogOpened",
    (current: CatalogViewerState) => {
        const albumId = getDisplayedAlbumId(current);

        if (!albumId) {
            return current;
        }

        const album = current.allAlbums.find(a => albumIdEquals(a.albumId, albumId));

        if (!album) {
            return current;
        }

        return {
            ...current,
            dialog: {
                type: "EditNameDialog",
                albumId,
                albumName: album.name,
                customFolderName: "",
                isCustomFolderNameEnabled: false,
                isLoading: false,
                error: {},
            },
        };
    }
);

export type EditNameDialogOpened = ReturnType<typeof editNameDialogOpened>;
