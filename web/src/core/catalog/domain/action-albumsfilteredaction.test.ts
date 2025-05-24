import { reduceAlbumsFiltered, albumsFilteredAction } from "./action-albumsfilteredaction";
import { loadedStateWithTwoAlbums, twoAlbums } from "./tests/test-helper-state";

describe("reduceAlbumsFiltered", () => {
    it("should show only directly owned album after the AlbumsFilteredAction", () => {
        const got = reduceAlbumsFiltered(
            loadedStateWithTwoAlbums,
            albumsFilteredAction({ selfOwned: true, owners: [] })
        );
        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[0],
            albums: [twoAlbums[0]],
        });
    });

    it("should show all albums when the filter moves back to 'All albums'", () => {
        const got = reduceAlbumsFiltered(
            {
                ...loadedStateWithTwoAlbums,
                albums: [],
            },
            albumsFilteredAction({ owners: [] })
        );
        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[1],
            albums: twoAlbums,
        });
    });

    it("should filter albums to those with a certain owner when the filter with that owner is selected", () => {
        const got = reduceAlbumsFiltered(
            {
                ...loadedStateWithTwoAlbums,
                albums: [],
            },
            albumsFilteredAction({ owners: ["herself"] })
        );
        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[2],
            albums: [twoAlbums[1]],
        });
    });
});
