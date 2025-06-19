import {initialCatalogState} from "../../language/initial-catalog-state";
import {myselfUser, summerTripAlbum, winterHolidaysAlbum} from "../tests/test-helper-state";
import {reduceEditAlbumDatesDialogOpened, editAlbumDatesDialogOpened} from "./action-editAlbumDatesDialogOpened";
import {editAlbumDatesDialogSelector} from "./selector-editAlbumDatesDialogSelector";
import {CatalogViewerState} from "../../language";

describe("action:editAlbumDatesDialogOpened", () => {
    const baseState = initialCatalogState(myselfUser);

    it("opens the dialog and sets the album ID for Summer Trip", () => {
        const stateWithAlbums: CatalogViewerState = {
            ...baseState,
            albums: [summerTripAlbum, winterHolidaysAlbum],
            allAlbums: [summerTripAlbum, winterHolidaysAlbum],
        };

        const action = editAlbumDatesDialogOpened(summerTripAlbum.albumId);
        const newState = reduceEditAlbumDatesDialogOpened(stateWithAlbums, action);

        const selection = editAlbumDatesDialogSelector(newState);

        expect(selection).toEqual({
            isOpen: true,
            albumName: "Summer Trip",
            startDate: new Date("2023-07-01T00:00:00.000Z"),
            endDate: new Date("2023-07-31T23:59:59.999Z"), // Inclusive end date
            isStartDateAtStartOfDay: true,
            isEndDateAtEndOfDay: true,
        });
    });

    it("opens the dialog and sets the album ID for Winter Holidays", () => {
        const stateWithAlbums: CatalogViewerState = {
            ...baseState,
            albums: [summerTripAlbum, winterHolidaysAlbum],
            allAlbums: [summerTripAlbum, winterHolidaysAlbum],
        };

        const action = editAlbumDatesDialogOpened(winterHolidaysAlbum.albumId);
        const newState = reduceEditAlbumDatesDialogOpened(stateWithAlbums, action);

        const selection = editAlbumDatesDialogSelector(newState);

        expect(selection).toEqual({
            isOpen: true,
            albumName: "Winter Holidays",
            startDate: new Date("2024-12-20T00:00:00.000Z"),
            endDate: new Date("2025-01-04T23:59:59.999Z"), // Inclusive end date
            isStartDateAtStartOfDay: true,
            isEndDateAtEndOfDay: true,
        });
    });

    it("should not change other state properties", () => {
        const stateWithAlbums: CatalogViewerState = {
            ...baseState,
            albums: [summerTripAlbum, winterHolidaysAlbum],
            allAlbums: [summerTripAlbum, winterHolidaysAlbum],
            mediasLoaded: true, // Example of another property
        };

        const action = editAlbumDatesDialogOpened(summerTripAlbum.albumId);
        const newState = reduceEditAlbumDatesDialogOpened(stateWithAlbums, action);

        expect(newState.mediasLoaded).toBe(true);
        expect(newState.shareModal).toBeUndefined();
        expect(newState.deleteDialog).toBeUndefined();
    });
});
