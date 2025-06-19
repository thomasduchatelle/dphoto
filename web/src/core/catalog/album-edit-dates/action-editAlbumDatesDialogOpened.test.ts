import {CatalogViewerState, initialCatalogState} from "../language";
import {myselfUser, twoAlbums} from "../tests/test-helper-state";
import {editAlbumDatesDialogOpened, reduceEditAlbumDatesDialogOpened} from "./action-editAlbumDatesDialogOpened";
import {editAlbumDatesDialogSelector} from "./selector-editAlbumDatesDialogSelector";

describe("action:editAlbumDatesDialogOpened", () => {
    const baseState = initialCatalogState(myselfUser);
    const jan2025Album = twoAlbums[0];
    const feb2025Album = twoAlbums[1];

    it("opens the dialog with an inclusive-exclusive dates range for the album ID for January 2025", () => {
        const stateWithAlbums: CatalogViewerState = {
            ...baseState,
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

    it("opens the dialog even when the AlbumId doesn't exist in the current albums list, displaying default values", () => {
        const nonExistentAlbumId = {owner: "unknown", folderName: "non-existent"};
        const stateWithoutAlbum: CatalogViewerState = {
            ...baseState,
            albums: [jan2025Album], // Only one album, not the one we're looking for
            allAlbums: [jan2025Album],
        };

        const action = editAlbumDatesDialogOpened(nonExistentAlbumId);
        const newState = reduceEditAlbumDatesDialogOpened(stateWithoutAlbum, action);

        const selection = editAlbumDatesDialogSelector(newState);

        expect(selection.isOpen).toBe(true);
        expect(selection.albumName).toBe(""); // Default empty string
        expect(selection.startDate).toEqual(expect.any(Date)); // Default new Date()
        expect(selection.endDate).toEqual(expect.any(Date));   // Default new Date()
        expect(selection.isStartDateAtStartOfDay).toBe(true);
        expect(selection.isEndDateAtEndOfDay).toBe(true);

        // Verify that the state itself reflects the dialog being open with the non-existent albumId
        expect(newState.editAlbumDatesDialog).toEqual({albumId: nonExistentAlbumId});
    });
});
