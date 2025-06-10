import {albumsFiltered, reduceAlbumsFiltered} from "./action-albumsFiltered";
import {herselfOwner, loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";

describe("action:albumsFiltered", () => {
    it("should show only directly owned album after the AlbumsFiltered", () => {
        const got = reduceAlbumsFiltered(
            loadedStateWithTwoAlbums,
            albumsFiltered({criterion: {selfOwned: true, owners: []}})
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
            albumsFiltered({criterion: {owners: []}})
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
            albumsFiltered({criterion: {owners: [herselfOwner]}})
        );
        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[2],
            albums: [twoAlbums[1]],
        });
    });
});
