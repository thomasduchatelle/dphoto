import {editDatesDialogOpened} from "./action-editDatesDialogOpened";
import {CatalogViewerState} from "../language";
import {loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {DEFAULT_EDIT_DATES_DIALOG_SELECTION, editDatesDialogSelector} from "./selector-editDatesDialogSelector";

const jan25Album = twoAlbums[0];
const feb25Album = twoAlbums[1];

describe("action:editDatesDialogOpened", () => {
    it("opens the dialog with the currently selected album data", () => {
        const action = editDatesDialogOpened();
        const state = action.reducer(loadedStateWithTwoAlbums, action);
        const got = editDatesDialogSelector(state);

        expect(got).toEqual({
            ...DEFAULT_EDIT_DATES_DIALOG_SELECTION,
            isOpen: true,
            albumName: jan25Album.name,
            startDate: new Date(2025, 0, 1),
            endDate: new Date(2025, 0, 31),
        });
    });

    it("opens the dialog with the album specified by loadingMediasFor", () => {
        const initialState: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            mediasLoadedFromAlbumId: jan25Album.albumId, // This should be ignored
            loadingMediasFor: feb25Album.albumId,
        };

        const action = editDatesDialogOpened();
        const state = action.reducer(initialState, action);
        const got = editDatesDialogSelector(state);

        expect(got).toEqual({
            ...DEFAULT_EDIT_DATES_DIALOG_SELECTION,
            isOpen: true,
            albumName: feb25Album.name,
            startDate: new Date(2025, 1, 1),
            endDate: new Date(2025, 2, 0),
        });
    });

    it("does not open dialog when no album is selected", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            albums: [jan25Album, feb25Album],
            mediasLoadedFromAlbumId: undefined,
            loadingMediasFor: undefined,
        };

        const action = editDatesDialogOpened();
        const got = action.reducer(state, action);
        const dialogSelection = editDatesDialogSelector(got);

        expect(dialogSelection.isOpen).toBeFalsy();
        expect(dialogSelection).toEqual(DEFAULT_EDIT_DATES_DIALOG_SELECTION);
        expect(got).toEqual(state); // Ensure no other state changes
    });

    it("does not open dialog when selected album is not found", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            albums: [jan25Album],
            mediasLoadedFromAlbumId: feb25Album.albumId,
        };

        const action = editDatesDialogOpened();
        const got = action.reducer(state, action);
        const dialogSelection = editDatesDialogSelector(got);

        expect(dialogSelection.isOpen).toBeFalsy();
        expect(dialogSelection).toEqual(DEFAULT_EDIT_DATES_DIALOG_SELECTION);
        expect(got).toEqual(state); // Ensure no other state changes
    });

    it("supports action comparison for testing", () => {
        const action1 = editDatesDialogOpened();
        const action2 = editDatesDialogOpened();
        
        expect(action1).toEqual(action2);
        expect([action1]).toContain(action2);
    });
});
