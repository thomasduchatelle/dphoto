import {createAction} from "@/libs/daction";
import {albumIdEquals, CatalogViewerState, editNameDialogNoError} from "../language";
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
                originalAlbumName: album.albumId.folderName,
                originalFolderName: album.albumId.folderName,
                customFolderName: "",
                isCustomFolderNameEnabled: false,
                isLoading: false,
                nameError: editNameDialogNoError,
            },
        };
    }
);

export type EditNameDialogOpened = ReturnType<typeof editNameDialogOpened>;
