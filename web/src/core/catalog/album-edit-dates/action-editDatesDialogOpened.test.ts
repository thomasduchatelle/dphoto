import {loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {CatalogViewerState} from "../language";
import {editDatesDialogOpenedAction, reduceEditDatesDialogOpened} from "./action-editDatesDialogOpened";
import {selectEditDatesDialog} from "./selector-editDatesDialog";

describe("action:editDatesDialogOpened", () => {
    it("returns isOpen: false and albumName: '' when no action is reduced", () => {
        const got = selectEditDatesDialog(loadedStateWithTwoAlbums);

        expect(got).toEqual({
            isOpen: false,
            albumName: ""
        });
    });

    it("returns isOpen: true and the correct album name when action is dispatched with ID of one of the albums", () => {
        const state = reduceEditDatesDialogOpened(
            loadedStateWithTwoAlbums,
            editDatesDialogOpenedAction(twoAlbums[0].albumId)
        );

        const got = selectEditDatesDialog(state);

        expect(got).toEqual({
            isOpen: true,
            albumName: "January 2025"
        });
    });

    it("returns isOpen: false and albumName: '' when action is dispatched with ID of a different album", () => {
        const nonExistentAlbumId = {owner: "unknown", folderName: "unknown"};
        
        const state = reduceEditDatesDialogOpened(
            loadedStateWithTwoAlbums,
            editDatesDialogOpenedAction(nonExistentAlbumId)
        );

        const got = selectEditDatesDialog(state);

        expect(got).toEqual({
            isOpen: false,
            albumName: ""
        });
    });

    it("returns isOpen: true and the second album name when dialog is open on first album and action is dispatched with ID of second album", () => {
        const stateWithDialogOpen: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: {
                albumId: twoAlbums[0].albumId
            }
        };

        const state = reduceEditDatesDialogOpened(
            stateWithDialogOpen,
            editDatesDialogOpenedAction(twoAlbums[1].albumId)
        );

        const got = selectEditDatesDialog(state);

        expect(got).toEqual({
            isOpen: true,
            albumName: "February 2025"
        });
    });
});
