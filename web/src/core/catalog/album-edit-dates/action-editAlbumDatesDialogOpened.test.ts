import {initialCatalogState} from "../../language/initial-catalog-state";
import {myselfUser, twoAlbums} from "../tests/test-helper-state";
import {reduceEditAlbumDatesDialogOpened, editAlbumDatesDialogOpened} from "./action-editAlbumDatesDialogOpened";
import {editAlbumDatesDialogSelector} from "./selector-editAlbumDatesDialogSelector";
import {CatalogViewerState} from "../../language";

describe("action:editAlbumDatesDialogOpened", () => {
    const baseState = initialCatalogState(myselfUser);
    const jan2025Album = twoAlbums[0];
    const feb2025Album = twoAlbums[1];

    it("opens the dialog and sets the album ID for January 2025", () => {
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

    it("opens the dialog and sets the album ID for February 2025", () => {
        const stateWithAlbums: CatalogViewerState = {
            ...baseState,
            albums: [jan2025Album, feb2025Album],
            allAlbums: [jan2025Album, feb2025Album],
        };

        const action = editAlbumDatesDialogOpened(feb2025Album.albumId);
        const newState = reduceEditAlbumDatesDialogOpened(stateWithAlbums, action);

        const selection = editAlbumDatesDialogSelector(newState);

        expect(selection).toEqual({
            isOpen: true,
            albumName: "February 2025",
            startDate: new Date(2025, 1, 1),
            endDate: new Date(2025, 1, 28, 23, 59, 59, 999), // Inclusive end date (Feb 2025 has 28 days)
            isStartDateAtStartOfDay: true,
            isEndDateAtEndOfDay: true,
        });
    });

    it("should not change other state properties", () => {
        const stateWithAlbums: CatalogViewerState = {
            ...baseState,
            albums: [jan2025Album, feb2025Album],
            allAlbums: [jan2025Album, feb2025Album],
            mediasLoaded: true, // Example of another property
        };

        const action = editAlbumDatesDialogOpened(jan2025Album.albumId);
        const newState = reduceEditAlbumDatesDialogOpened(stateWithAlbums, action);

        expect(newState.mediasLoaded).toBe(true);
        expect(newState.shareModal).toBeUndefined();
        expect(newState.deleteDialog).toBeUndefined();
    });
});
