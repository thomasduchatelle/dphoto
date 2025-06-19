import {CatalogViewerState} from "../language";
import {loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {editAlbumDatesDialogClosed, reduceEditAlbumDatesDialogClosed} from "./action-editAlbumDatesDialogClosed";
import {editAlbumDatesDialogSelector} from "./selector-editAlbumDatesDialogSelector";

describe("action:editAlbumDatesDialogClosed", () => {
    const jan2025Album = twoAlbums[0];

    const editDateDialogOpenForAlbum0: Partial<CatalogViewerState> = {
        editAlbumDatesDialog: {
            albumId: jan2025Album.albumId,
        },
    }

    it("closes the dialog when receiving the action editAlbumDatesDialogClosed", () => {
        const stateWithDialogOpened: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            ...editDateDialogOpenForAlbum0,
        };

        const action = editAlbumDatesDialogClosed();
        const newState = reduceEditAlbumDatesDialogClosed(stateWithDialogOpened, action);

        const selection = editAlbumDatesDialogSelector(newState);

        expect(selection.isOpen).toEqual(false);
    });

    it("should handle closing when dialog is already closed", () => {
        const action = editAlbumDatesDialogClosed();
        const newState = reduceEditAlbumDatesDialogClosed(loadedStateWithTwoAlbums, action);

        const selection = editAlbumDatesDialogSelector(newState);

        expect(selection.isOpen).toBe(false);
    });
});
