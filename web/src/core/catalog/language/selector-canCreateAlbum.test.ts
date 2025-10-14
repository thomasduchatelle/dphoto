import {canCreateAlbumSelector, CanCreateAlbumSelection} from "./selector-canCreateAlbum";
import {Album, CatalogViewerState} from "./catalog-state";
import {loadedStateWithTwoAlbums, myselfUser, herselfUser, herselfOwner} from "../tests/test-helper-state";

const visitorAlbum: Album = {
    albumId: {owner: herselfOwner, folderName: "visitor-album"},
    name: "Visitor Album",
    start: new Date(2025, 0, 1),
    end: new Date(2025, 1, 1),
    totalCount: 10,
    temperature: 0.25,
    relativeTemperature: 1,
    ownedBy: {name: "Herself", users: [herselfUser]},
    sharedWith: [],
};

const stateAsVisitor: CatalogViewerState = {
    ...loadedStateWithTwoAlbums,
    allAlbums: [visitorAlbum],
};

describe('selector:canCreateAlbumSelector', () => {
    it('should return canCreateAlbum=true when user has owned albums', () => {
        const got = canCreateAlbumSelector(loadedStateWithTwoAlbums);

        expect(got).toEqual<CanCreateAlbumSelection>({
            canCreateAlbum: true,
        });
    });

    it('should return canCreateAlbum=false when user has no owned albums (visitor)', () => {
        const got = canCreateAlbumSelector(stateAsVisitor);

        expect(got).toEqual<CanCreateAlbumSelection>({
            canCreateAlbum: false,
        });
    });

    it('should return canCreateAlbum=false when there are no albums at all', () => {
        const emptyState: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            allAlbums: [],
        };

        const got = canCreateAlbumSelector(emptyState);

        expect(got).toEqual<CanCreateAlbumSelection>({
            canCreateAlbum: false,
        });
    });
});
