import {initialCatalogState} from "../../language/initial-catalog-state";
import {myselfUser, summerTripAlbum} from "../tests/test-helper-state";
import {reduceEditAlbumDatesDialogClosed, editAlbumDatesDialogClosed} from "./action-editAlbumDatesDialogClosed";
import {editAlbumDatesDialogSelector} from "./selector-editAlbumDatesDialogSelector";
import {CatalogViewerState} from "../../language";

describe("action:editAlbumDatesDialogClosed", () => {
    const baseState = initialCatalogState(myselfUser);

    it("closes the dialog and clears the album ID", () => {
        const stateWithDialogOpened: CatalogViewerState = {
            ...baseState,
            editAlbumDatesDialog: {
                albumId: summerTripAlbum.albumId,
            },
            albums: [summerTripAlbum], // Ensure album is present for selector to find it
            allAlbums: [summerTripAlbum],
        };

        const action = editAlbumDatesDialogClosed();
        const newState = reduceEditAlbumDatesDialogClosed(stateWithDialogOpened, action);

        const selection = editAlbumDatesDialogSelector(newState);

        expect(selection).toEqual({
            isOpen: false,
            albumName: "",
            startDate: expect.any(Date), // Default date when dialog is closed
            endDate: expect.any(Date),   // Default date when dialog is closed
            isStartDateAtStartOfDay: true,
            isEndDateAtEndOfDay: true,
        });
        expect(newState.editAlbumDatesDialog).toBeUndefined();
    });

    it("should not change other state properties", () => {
        const stateWithDialogOpened: CatalogViewerState = {
            ...baseState,
            editAlbumDatesDialog: {
                albumId: summerTripAlbum.albumId,
            },
            albums: [summerTripAlbum],
            allAlbums: [summerTripAlbum],
            mediasLoaded: true, // Example of another property
        };

        const action = editAlbumDatesDialogClosed();
        const newState = reduceEditAlbumDatesDialogClosed(stateWithDialogOpened, action);

        expect(newState.mediasLoaded).toBe(true);
        expect(newState.shareModal).toBeUndefined();
        expect(newState.deleteDialog).toBeUndefined();
    });

    it("should handle closing when dialog is already closed", () => {
        const action = editAlbumDatesDialogClosed();
        const newState = reduceEditAlbumDatesDialogClosed(baseState, action);

        const selection = editAlbumDatesDialogSelector(newState);

        expect(selection.isOpen).toBe(false);
        expect(newState.editAlbumDatesDialog).toBeUndefined();
    });
});
