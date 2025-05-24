import { CatalogViewerState, Album, MediaWithinADay, MediaType, SharingType, UserDetails } from "./catalog-state";
import { CurrentUserInsight } from "./catalog-reducer";
import { AlbumsAndMediasLoadedAction, makeReduceAlbumsAndMediasLoaded } from "./catalog-action-AlbumsAndMediasLoadedAction";

describe("reduceAlbumsAndMediasLoaded", () => {
    const myselfUser: CurrentUserInsight = {picture: "my-face.jpg"};
    const herselfUser: UserDetails = {email: "her@self.com", name: "Herself", picture: "her-face.jpg"};
    const herselfOwner = "herself";

    const twoAlbums: Album[] = [
        {
            albumId: {owner: "myself", folderName: "jan-25"},
            name: "January 2025",
            start: new Date(2025, 0, 1),
            end: new Date(2025, 0, 31),
            totalCount: 42,
            temperature: 0.25,
            relativeTemperature: 1,
            sharedWith: [
                {
                    user: herselfUser,
                    role: SharingType.visitor,
                }
            ],
        },
        {
            albumId: {owner: herselfOwner, folderName: "feb-25"},
            name: "February 2025",
            start: new Date(2025, 1, 1),
            end: new Date(2025, 2, 0),
            totalCount: 12,
            temperature: 0.25,
            relativeTemperature: 1,
            ownedBy: {name: "Herself", users: [herselfUser]},
            sharedWith: [],
        },
    ];

    const someMedias: MediaWithinADay[] = [{
        day: new Date(2025, 0, 1), medias: [{
            id: "media-1",
            type: MediaType.IMAGE,
            time: new Date("2025-01-05T12:42:00Z"),
            uiRelativePath: "media-1/image.jpg",
            contentPath: "/media-1.jpg",
            source: "",
        }]
    }];

    const initialCatalogState: CatalogViewerState = {
        albumFilterOptions: [],
        albumFilter: undefined as any,
        allAlbums: [],
        albumNotFound: false,
        albums: [],
        medias: [],
        albumsLoaded: false,
        mediasLoaded: false
    };

    it("should add the loaded albums and medias to the state, and reset all status", () => {
        const reduce = makeReduceAlbumsAndMediasLoaded(myselfUser);
        const action = AlbumsAndMediasLoadedAction(twoAlbums, someMedias, twoAlbums[0]);
        const got = reduce({
            ...initialCatalogState,
            albumNotFound: true,
            albumsLoaded: false,
            mediasLoaded: false,
        }, action);

        expect(got.allAlbums).toEqual(twoAlbums);
        expect(got.albums).toEqual(twoAlbums);
        expect(got.medias).toEqual(someMedias);
        expect(got.albumNotFound).toBe(false);
        expect(got.albumsLoaded).toBe(true);
        expect(got.mediasLoaded).toBe(true);
        expect(got.mediasLoadedFromAlbumId).toEqual(twoAlbums[0].albumId);
        expect(got.albumFilterOptions.length).toBeGreaterThan(0);
        expect(got.albumFilter).toBeDefined();
    });

    it("should use 'All albums' filter even when it's the only selection available (only directly owned albums)", () => {
        const reduce = makeReduceAlbumsAndMediasLoaded(myselfUser);
        const action = AlbumsAndMediasLoadedAction([twoAlbums[0]], someMedias, twoAlbums[0]);
        const got = reduce(initialCatalogState, action);

        expect(got.allAlbums).toEqual([twoAlbums[0]]);
        expect(got.albums).toEqual([twoAlbums[0]]);
        expect(got.albumFilterOptions.length).toBe(1);
        expect(got.albumFilter.name).toBe("All albums");
    });
});
