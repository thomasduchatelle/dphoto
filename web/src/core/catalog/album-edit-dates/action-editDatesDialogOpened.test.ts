import {editDatesDialogOpened, reduceEditDatesDialogOpened} from "./action-editDatesDialogOpened";
import {CatalogViewerState} from "../language";
import {loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {editDatesDialogSelector} from "./selector-editDatesDialogSelector";

const jan25Album = twoAlbums[0];
const feb25Album = twoAlbums[1];

describe("action:editDatesDialogOpened", () => {
    it("opens the dialog with the currently selected album data", () => {
        const state = reduceEditDatesDialogOpened(loadedStateWithTwoAlbums, editDatesDialogOpened());
        const got = editDatesDialogSelector(state);

        expect(got).toEqual({
            isOpen: true,
            albumName: jan25Album.name,
            startDate: new Date(2025, 0, 1),
            endDate: new Date(2025, 0, 31),
            startAtDayStart: true,
            endAtDayEnd: true,
            isLoading: false,
        });
    });

    it("opens the dialog with the album specified by loadingMediasFor", () => {
        const initialState: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            mediasLoadedFromAlbumId: jan25Album.albumId, // This should be ignored
            loadingMediasFor: feb25Album.albumId,
        };

        const state = reduceEditDatesDialogOpened(initialState, editDatesDialogOpened());
        const got = editDatesDialogSelector(state);

        expect(got).toEqual({
            isOpen: true,
            albumName: feb25Album.name,
            startDate: new Date(2025, 1, 1),
            endDate: new Date(2025, 2, 0),
            startAtDayStart: true,
            endAtDayEnd: true,
            isLoading: false,
        });
    });

    it("does not open dialog when no album is selected", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            albums: [jan25Album, feb25Album],
            mediasLoadedFromAlbumId: undefined,
            loadingMediasFor: undefined,
        };

        const got = reduceEditDatesDialogOpened(state, editDatesDialogOpened());
        const dialogSelection = editDatesDialogSelector(got);

        expect(dialogSelection.isOpen).toBeFalsy();
        expect(got).toEqual(state); // Ensure no other state changes
    });

    it("does not open dialog when selected album is not found", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            albums: [jan25Album],
            mediasLoadedFromAlbumId: feb25Album.albumId,
        };

        const got = reduceEditDatesDialogOpened(state, editDatesDialogOpened());
        const dialogSelection = editDatesDialogSelector(got);

        expect(dialogSelection.isOpen).toBeFalsy();
        expect(got).toEqual(state); // Ensure no other state changes
    });
});
