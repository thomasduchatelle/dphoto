import {Album, CatalogViewerState, OwnerDetails} from "../catalog-state";
import {openDeleteAlbumDialogAction, reduceOpenDeleteAlbumDialog} from "./delete-openDeleteAlbumDialog";
import {selectDeleteAlbumDialog} from "./delete-album-dialog-selector";
import {myselfUser, someMedias} from "../tests/test-helper-state";
import {generateAlbumFilterOptions} from "../catalog-common-modifiers";

const ownerDetails: OwnerDetails = {
    name: "Other Owner",
    users: [{ name: "Other", email: "other@example.com" }]
};

const albumSpring: Album = {
    albumId: {owner: "user1", folderName: "spring"},
    name: "Spring Album",
    start: new Date("2023-03-01"),
    end: new Date("2023-03-31"),
    totalCount: 10,
    temperature: 1,
    relativeTemperature: 1,
    sharedWith: [],
};

const albumSummer: Album = {
    albumId: {owner: "user1", folderName: "summer"},
    name: "Summer Album",
    start: new Date("2023-06-01"),
    end: new Date("2023-06-30"),
    totalCount: 8,
    temperature: 1,
    relativeTemperature: 1,
    sharedWith: [],
    ownedBy: ownerDetails, // Not deletable
};

const albumWinter: Album = {
    albumId: {owner: "user1", folderName: "winter"},
    name: "Winter Album",
    start: new Date("2023-12-01"),
    end: new Date("2023-12-31"),
    totalCount: 5,
    temperature: 1,
    relativeTemperature: 1,
    sharedWith: [],
};

const allAlbums = [albumSpring, albumSummer, albumWinter];

const albumFilterOptions = generateAlbumFilterOptions(myselfUser, allAlbums);
const stateWithThreeAlbumsLoaded: CatalogViewerState = {
    currentUser: myselfUser,
    allAlbums: [albumSpring, albumSummer, albumWinter],
    albumFilterOptions: albumFilterOptions,
    albumFilter: albumFilterOptions[0],
    albums: allAlbums,
    medias: someMedias,
    albumNotFound: false,
    mediasLoadedFromAlbumId: allAlbums[0].albumId,
    albumsLoaded: true,
    mediasLoaded: true,
};

describe("OpenDeleteAlbumDialog", () => {
    it("results to dialog open, pre-selected album to the one currently loaded, and albums only containing those deletable, when the dialog was not open", () => {
        const next = reduceOpenDeleteAlbumDialog(stateWithThreeAlbumsLoaded, openDeleteAlbumDialogAction());
        expect(selectDeleteAlbumDialog(next)).toEqual({
            isOpen: true,
            albums: [albumSpring, albumWinter],
            initialSelectedAlbumId: albumSpring.albumId,
            isLoading: false,
        });
    });

    it("results to dialog open, pre-selected album to the one being loaded, and albums only containing those deletable, when the dialog was not open and medias are currently being loaded", () => {
        const next = reduceOpenDeleteAlbumDialog({
            ...stateWithThreeAlbumsLoaded,
            loadingMediasFor: albumWinter.albumId,
        }, openDeleteAlbumDialogAction());

        expect(selectDeleteAlbumDialog(next)).toEqual({
            isOpen: true,
            albums: [albumSpring, albumWinter],
            initialSelectedAlbumId: albumWinter.albumId,
            isLoading: false,
        });
    });

    it("results to dialog open with first deletable album selected when a non-deletable album is the one loaded", () => {
        // mediasLoadedFromAlbumId is set to albumSummer (not deletable)
        const next = reduceOpenDeleteAlbumDialog({
            ...stateWithThreeAlbumsLoaded,
            mediasLoadedFromAlbumId: albumSummer.albumId,
        }, openDeleteAlbumDialogAction());

        expect(selectDeleteAlbumDialog(next)).toEqual({
            isOpen: true,
            albums: [albumSpring, albumWinter],
            initialSelectedAlbumId: albumSpring.albumId,
            isLoading: false,
        });
    });

    it("results to dialog open with deletable albums, no error and loading false, when the dialog were already open", () => {
        const next = reduceOpenDeleteAlbumDialog({
            ...stateWithThreeAlbumsLoaded,
            deleteDialog: {
                deletableAlbums: [albumWinter],
                initialSelectedAlbumId: albumWinter.albumId,
                isLoading: true,
                error: "Some error",
            },
        }, openDeleteAlbumDialogAction());

        expect(selectDeleteAlbumDialog(next)).toEqual({
            isOpen: true,
            albums: [albumSpring, albumWinter],
            initialSelectedAlbumId: albumSpring.albumId,
            error: undefined,
            isLoading: false,
        });
    });
});
