import {albumDeletedAction, reduceAlbumDeleted} from "./delete-albumDeletedAction";
import {loadedStateWithTwoAlbums, myselfUser, twoAlbums, twoAlbumsNoFilterOptions} from "../tests/test-helper-state";
import {initialCatalogState} from "../initial-catalog-state";
import {CatalogViewerState} from "../catalog-state";

describe("reduceAlbumDeleted", () => {
    const deleteDialog = {deletableAlbums: twoAlbums, isLoading: true};

    const marchAlbum = {
        albumId: {owner: "myself", folderName: "mar-25"},
        name: "March 25",
        start: new Date(2025, 2, 1),
        end: new Date(2025, 2, 31),
        totalCount: 0,
        temperature: 0,
        relativeTemperature: 0,
        sharedWith: []
    };

    const loadedStateWithThreeAlbums: CatalogViewerState = {
        ...loadedStateWithTwoAlbums,
        allAlbums: [...twoAlbums, marchAlbum],
        albums: [...twoAlbums, marchAlbum],
    }

    it("closes the dialog and update the lists of albums list like an initial loading", () => {
        const got = reduceAlbumDeleted(
            {
                ...initialCatalogState(myselfUser),
                deleteDialog,
            },
            albumDeletedAction({albums: twoAlbums})
        );

        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            medias: [],
            mediasLoaded: false,
            mediasLoadedFromAlbumId: undefined,
        });
    });

    it("closes the dialog and list all albums (with single filter option) when the only owned album has been removed", () => {
        const got = reduceAlbumDeleted(
            {
                ...loadedStateWithTwoAlbums,
                albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[1], // only owned album
                deleteDialog,
            },
            albumDeletedAction({albums: [twoAlbums[1]], redirectTo: twoAlbums[1].albumId})
        );

        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            allAlbums: [twoAlbums[1]],
            albums: [twoAlbums[1]],
            albumFilterOptions: [
                twoAlbumsNoFilterOptions,
            ],
            albumFilter: twoAlbumsNoFilterOptions,
        });
    });

    it("closes the dialog and hold the filter when the filtered album list is not empty and there is no redirect", () => {
        const got = reduceAlbumDeleted(
            {
                ...loadedStateWithThreeAlbums,
                albumFilter: loadedStateWithThreeAlbums.albumFilterOptions[0],
                albums: [loadedStateWithThreeAlbums.allAlbums[0], loadedStateWithThreeAlbums.allAlbums[2]],
                deleteDialog,
            },
            albumDeletedAction({albums: twoAlbums})
        );
        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            albums: [loadedStateWithThreeAlbums.allAlbums[0]],
            albumFilter: loadedStateWithThreeAlbums.albumFilterOptions[0],
        });
    });

    it("closes the dialog and hold the filter when the filtered album list is not empty and contains the albumId redirected to", () => {
        const got = reduceAlbumDeleted(
            {
                ...loadedStateWithThreeAlbums,
                albumFilter: loadedStateWithThreeAlbums.albumFilterOptions[0],
                albums: [loadedStateWithThreeAlbums.allAlbums[0], loadedStateWithThreeAlbums.allAlbums[2]],
                deleteDialog,
            },
            albumDeletedAction({albums: twoAlbums, redirectTo: twoAlbums[0].albumId})
        );
        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            albums: [loadedStateWithThreeAlbums.allAlbums[0]],
            albumFilter: loadedStateWithThreeAlbums.albumFilterOptions[0],
        });
    });


    it('closes the dialog and changes the filter to "All Albums" when the original filter would filter out the selected album', () => {
        const got = reduceAlbumDeleted(
            {
                ...loadedStateWithThreeAlbums,
                albumFilter: loadedStateWithThreeAlbums.albumFilterOptions[0],
                albums: [loadedStateWithThreeAlbums.allAlbums[0], loadedStateWithThreeAlbums.allAlbums[2]],
                deleteDialog,
            },
            albumDeletedAction({albums: twoAlbums, redirectTo: twoAlbums[1].albumId})
        );
        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
        });
    });
});
