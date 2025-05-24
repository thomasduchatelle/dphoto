import { CatalogViewerState, Album, MediaWithinADay, MediaType, SharingType, UserDetails } from "./catalog-state";
import { AlbumsLoadedAction, reduceAlbumsLoaded } from "./catalog-action-AlbumsLoadedAction";

describe("reduceAlbumsLoaded", () => {
    const myselfUser = {picture: "my-face.jpg"};
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
        currentUser: myselfUser,
        albumFilterOptions: [],
        albumFilter: undefined as any,
        allAlbums: [],
        albumNotFound: false,
        albums: [],
        medias: [],
        albumsLoaded: false,
        mediasLoaded: false
    };

    it("should update the list of albums and clear errors when AlbumsLoadedAction is received", () => {
        const action = AlbumsLoadedAction(twoAlbums);
        const got = reduceAlbumsLoaded({
            ...initialCatalogState,
            allAlbums: [twoAlbums[0]],
            albums: [twoAlbums[0]],
            error: new Error("TEST previous error to clear"),
            albumsLoaded: false,
        }, action);

        expect(got.allAlbums).toEqual(twoAlbums);
        expect(got.albums).toEqual(twoAlbums);
        expect(got.error).toBeUndefined();
        expect(got.albumsLoaded).toBe(true);
    });

    it("should update the available filters and re-apply the selected filter when receiving AlbumsLoadedAction", () => {
        const action = AlbumsLoadedAction(twoAlbums, twoAlbums[0].albumId);
        const got = reduceAlbumsLoaded({
            ...initialCatalogState,
            albumFilterOptions: [],
            albumFilter: undefined as any,
            allAlbums: [twoAlbums[1]],
            albums: [],
            mediasLoadedFromAlbumId: twoAlbums[0].albumId, // no effect
        }, action);

        expect(got.allAlbums).toEqual(twoAlbums);
        expect(got.albums).toEqual([twoAlbums[0]]);
        expect(got.albumFilter).toBeDefined();
    });

    it("should remove the album filter if the redirectTo in AlbumsLoadedAction wouldn't be displayed", () => {
        const action = AlbumsLoadedAction(twoAlbums, twoAlbums[1].albumId);
        const got = reduceAlbumsLoaded({
            ...initialCatalogState,
            allAlbums: [twoAlbums[0]],
            albums: [],
            albumFilter: undefined as any,
        }, action);

        expect(got.allAlbums).toEqual(twoAlbums);
        expect(got.albums).toEqual(twoAlbums);
        expect(got.albumFilter).toBeDefined();
    });
});
