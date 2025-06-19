import {CatalogViewerState} from "../language";
import {loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {editAlbumDatesDialogOpened, reduceEditAlbumDatesDialogOpened} from "./action-editAlbumDatesDialogOpened";
import {editAlbumDatesDialogSelector} from "./selector-editAlbumDatesDialogSelector";

describe("action:editAlbumDatesDialogOpened", () => {
    const jan2025Album = twoAlbums[0];
    const feb2025Album = twoAlbums[1];

    it("opens the dialog with an inclusive-exclusive dates range for the album ID for January 2025", () => {
        const stateWithAlbums: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            albums: [jan2025Album, feb2025Album],
            allAlbums: [jan2025Album, feb2025Album],
        };

        const action = editAlbumDatesDialogOpened(jan2025Album.albumId);
        const newState = reduceEditAlbumDatesDialogOpened(stateWithAlbums, action);

        const selection = editAlbumDatesDialogSelector(newState);

        expect(selection).toEqual({
            isOpen: true,
            albumName: "January 2025",
            startDate: new Date(2025, 0, 1),
            endDate: new Date(2025, 0, 31, 23, 59, 59, 999), // Inclusive end date
            isStartDateAtStartOfDay: true,
            isEndDateAtEndOfDay: true,
        });
    });

    it("should not open the dialog when the AlbumId doesn't exist", () => {
        const nonExistentAlbumId = {owner: "unknown", folderName: "non-existent"};
        const stateWithoutAlbum: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            albums: [jan2025Album], // Only one album, not the one we're looking for
            allAlbums: [jan2025Album],
        };

        const action = editAlbumDatesDialogOpened(nonExistentAlbumId);
        const newState = reduceEditAlbumDatesDialogOpened(stateWithoutAlbum, action);

        const selection = editAlbumDatesDialogSelector(newState);

        expect(selection.isOpen).toBe(false);
    });
});
