import {editDatesDialogOpened, reduceEditDatesDialogOpened} from "./action-editDatesDialogOpened";
import {CatalogViewerState} from "../language";
import {loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";

const jan25Album = twoAlbums[0];
const feb25Album = twoAlbums[1];

describe("action:editDatesDialogOpened", () => {
    it("opens the dialog with the currently selected album data", () => {
        const got = reduceEditDatesDialogOpened(loadedStateWithTwoAlbums, editDatesDialogOpened());

        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            editDatesDialog: {
                albumId: jan25Album.albumId,
                albumName: jan25Album.name,
                startDate: new Date(2025, 0, 1),
                endDate: new Date(2025, 0, 31),
            },
        });
    });

    it("opens the dialog with the album specified by loadingMediasFor", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            mediasLoadedFromAlbumId: jan25Album.albumId, // This should be ignored
            loadingMediasFor: feb25Album.albumId,
        };

        const got = reduceEditDatesDialogOpened(state, editDatesDialogOpened());

        expect(got).toEqual({
            ...state,
            editDatesDialog: {
                albumId: feb25Album.albumId,
                albumName: feb25Album.name,
                startDate: new Date(2025, 1, 1),
                endDate: new Date(2025, 2, 0),
            },
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

        expect(got).toEqual(state);
    });

    it("does not open dialog when selected album is not found", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            albums: [jan25Album],
            mediasLoadedFromAlbumId: feb25Album.albumId,
        };

        const got = reduceEditDatesDialogOpened(state, editDatesDialogOpened());

        expect(got).toEqual(state);
    });
});
