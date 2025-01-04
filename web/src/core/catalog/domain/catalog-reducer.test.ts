import {catalogReducerFunction, CurrentUserInsight, generateAlbumFilterOptions, initialCatalogState} from "./catalog-reducer";
import {Album, AlbumFilterEntry, CatalogViewerState, MediaType, MediaWithinADay, UserDetails} from "./catalog-state";

describe("CatalogViewerState", () => {
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
            // ownedBy: {name: "Myself", users: [myselfUser]}, TODO make possible to have details about the current user
            sharedWith: [],
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
    ]

    const twoAlbumsNoFilterOptions: AlbumFilterEntry = {
        criterion: {
            owners: []
        },
        avatars: [`${myselfUser.picture}`, `${herselfUser.picture}`],
        name: "All albums",
    };

    const someMedias: MediaWithinADay[] = [{
        day: new Date(2025, 0, 1), medias: [{
            id: "media-1",
            type: MediaType.IMAGE,
            time: new Date("2025-01-05T12:42:00Z"),
            uiRelativePath: "media-1/image.jpg",
            contentPath: "/media-1.jpg",
            source: "",
        }]
    }]

    const loadedStateWithTwoAlbums: CatalogViewerState = {
        allAlbums: twoAlbums,
        albumFilterOptions: [
            {
                criterion: {
                    selfOwned: true,
                    owners: [],
                },
                avatars: [myselfUser.picture ?? ""],
                name: "My albums",
            },
            twoAlbumsNoFilterOptions,
            {
                criterion: {
                    owners: [herselfOwner]
                },
                avatars: [herselfUser.picture ?? ""],
                name: herselfUser.name,
            },
        ],
        albumFilter: twoAlbumsNoFilterOptions,
        albums: twoAlbums,
        medias: someMedias,
        albumNotFound: false,
        albumsLoaded: true,
        mediasLoaded: true,
    };

    it("should add the loaded albums and medias to the state, and reset all status when receiving AlbumsAndMediasLoadedAction", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        const got = catalogReducer({
            ...initialCatalogState,
            albumNotFound: true,
            albumsLoaded: false,
            mediasLoaded: false,
        }, {
            type: "AlbumsAndMediasLoadedAction",
            albums: twoAlbums,
            medias: someMedias,
            selectedAlbum: twoAlbums[0],
        });

        expect(got).toEqual(loadedStateWithTwoAlbums)
    })

    it("should use 'All albums' filter even when it's the only selection available (only directly owned albums) when receiving AlbumsAndMediasLoadedAction", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        const got = catalogReducer(initialCatalogState, {
            type: "AlbumsAndMediasLoadedAction",
            albums: [twoAlbums[0]],
            medias: someMedias,
            selectedAlbum: twoAlbums[0],
        });

        const allAlbumFilter = {
            criterion: {
                owners: []
            },
            avatars: [myselfUser.picture ?? ""],
            name: "All albums",
        };
        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            albums: [twoAlbums[0]],
            allAlbums: [twoAlbums[0]],
            albumFilter: allAlbumFilter,
            albumFilterOptions: [allAlbumFilter],
        })
    })

    it("should show only directly owned album after the AlbumsFilteredAction", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        const got = catalogReducer(loadedStateWithTwoAlbums, {
            type: "AlbumsFilteredAction",
            criterion: {selfOwned: true, owners: []},
        });

        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[0],
            albums: [twoAlbums[0]],
        })
    })

    it("should show all albums when the filter moves back to 'All albums'", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        const got = catalogReducer({
            ...loadedStateWithTwoAlbums,
            albums: [],
        }, {
            type: "AlbumsFilteredAction",
            criterion: {owners: []},
        });

        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[1],
            albums: twoAlbums,
        })
    })

    it("should filter albums to those with a certain owner when the filter with that owner is selected", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        const got = catalogReducer({
            ...loadedStateWithTwoAlbums,
            albums: [],
        }, {
            type: "AlbumsFilteredAction",
            criterion: {owners: [herselfOwner]},
        });

        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[2],
            albums: [twoAlbums[1]],
        })
    })

    it("should only change the medias and loading status when reducing StartLoadingMediasAction, and clear errors", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        expect(catalogReducer({
            ...loadedStateWithTwoAlbums,
            albumNotFound: true,
            error: new Error("TEST previous error to clear"),
        }, {
            type: "StartLoadingMediasAction",
            albumId: twoAlbums[1].albumId,
        })).toEqual({
            ...loadedStateWithTwoAlbums,
            medias: [],
            loadingMediasFor: twoAlbums[1].albumId,
            mediasLoaded: false,
        })
    })

    it("should only change the medias and loading status when reducing MediasLoadedAction, and clear errors", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        expect(catalogReducer({
            ...loadedStateWithTwoAlbums,
            medias: [],
            mediasLoaded: false,
            loadingMediasFor: twoAlbums[1].albumId,
            albumNotFound: true,
            error: new Error("TEST previous error to clear"),
        }, {
            type: "MediasLoadedAction",
            albumId: twoAlbums[1].albumId,
            medias: someMedias,
        })).toEqual({
            ...loadedStateWithTwoAlbums,
            medias: someMedias,
        })
    })

    it("should ignore MediasLoadedAction if the medias are not for the expected album", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        expect(catalogReducer({
            ...loadedStateWithTwoAlbums,
            medias: [],
            loadingMediasFor: twoAlbums[0].albumId,
        }, {
            type: "MediasLoadedAction",
            albumId: twoAlbums[1].albumId,
            medias: someMedias,
        })).toEqual({
            ...loadedStateWithTwoAlbums,
            medias: [],
            loadingMediasFor: twoAlbums[0].albumId,
        })
    })

    it("should set the errors and clears medias and media loading status when reducing MediaFailedToLoadAction", () => {
        const testError = new Error("TEST loading error");

        const catalogReducer = catalogReducerFunction(myselfUser);
        expect(catalogReducer(initialCatalogState, {
            type: "MediaFailedToLoadAction",
            albums: twoAlbums,
            selectedAlbum: twoAlbums[0],
            error: testError,
        })).toEqual({
            ...loadedStateWithTwoAlbums,
            medias: [],
            error: testError,
        })
    })
})

describe("generateAlbumFilterOptions", () => {
    const currentUser = {picture: "my-face.jpg"};
    const herselfOwner = "herself";
    const herselfUser: UserDetails = {email: "her@self.com", name: "Herself", picture: "her-face.jpg"};
    const himselfOwner = "himself";
    const himself = {name: "Z-Himself", users: [{email: "him@self.com", name: "-", picture: "a-his-face.jpg"}]};
    const janBaseAlbum = {
        name: "January 2025",
        start: new Date(2025, 0, 1),
        end: new Date(2025, 1, 0),
        totalCount: 42,
        temperature: 0.25,
        relativeTemperature: 1,
        sharedWith: [],
    }
    const febBaseAlbum = {
        name: "February 2025",
        start: new Date(2025, 1, 1),
        end: new Date(2025, 2, 0),
        totalCount: 12,
        temperature: 0.25,
        relativeTemperature: 1,
        sharedWith: [],
    }

    it("should return the single option 'All albums' if no albums are available", () => {
        expect(generateAlbumFilterOptions(currentUser, [])).toEqual([{
            name: "All albums",
            criterion: {
                owners: []
            },
            avatars: [currentUser.picture],
        }])
    })

    it("should return the single option 'All albums' without picture if the current user doesn't have any", () => {
        expect(generateAlbumFilterOptions({}, [])).toEqual([{
            name: "All albums",
            criterion: {
                owners: []
            },
            avatars: [],
        }])
    })

    it("should return the single option 'All albums' if all albums are owned by the current user", () => {
        let janBaseAlbum = {
            name: "January 2025",
            start: new Date(2025, 0, 1),
            end: new Date(2025, 0, 31),
            totalCount: 42,
            temperature: 0.25,
            relativeTemperature: 1,
            sharedWith: [],
        };
        expect(generateAlbumFilterOptions(currentUser, [
            {
                ...janBaseAlbum,
                albumId: {owner: "myself", folderName: "jan-25"},
            },
        ])).toEqual([{
            name: "All albums",
            criterion: {
                owners: []
            },
            avatars: [currentUser.picture],
        }])
    })

    it("should return the single option 'All albums' if all albums are owned by another user", () => {
        expect(generateAlbumFilterOptions(currentUser, [
            {
                ...janBaseAlbum,
                albumId: {owner: "herself", folderName: "jan-25"},
                ownedBy: {name: "herself", users: [herselfUser]},
            },
        ])).toEqual([{
            name: "All albums",
            criterion: {
                owners: []
            },
            avatars: [currentUser.picture, herselfUser.picture],
        }])
    })

    it("should return the options 'My album', 'All albums', 'Her albums' if there is albums owned by current user and HER", () => {
        expect(generateAlbumFilterOptions(currentUser, [
            {
                ...janBaseAlbum,
                albumId: {owner: "myself", folderName: "jan-25"},
            },
            {
                ...febBaseAlbum,
                albumId: {owner: herselfOwner, folderName: "feb-25"},
                ownedBy: {name: herselfUser.name, users: [herselfUser]},
            },
        ])).toEqual([
            {
                name: "My albums",
                criterion: {
                    selfOwned: true,
                    owners: [],
                },
                avatars: [currentUser.picture],
            },
            {
                name: "All albums",
                criterion: {
                    owners: []
                },
                avatars: [currentUser.picture, herselfUser.picture],
            },
            {
                name: herselfUser.name,
                criterion: {
                    owners: [herselfOwner]
                },
                avatars: [herselfUser.picture],
            },
        ])
    })

    it("should return the options 'My album', 'All albums', 'Her albums' if there is albums owned by current user and HER even if HER doesn't have a picture", () => {
        expect(generateAlbumFilterOptions(currentUser, [
            {
                ...janBaseAlbum,
                albumId: {owner: "myself", folderName: "jan-25"},
            },
            {
                ...febBaseAlbum,
                albumId: {owner: "herself", folderName: "feb-25"},
                ownedBy: {name: herselfUser.name, users: [{...herselfUser, picture: undefined}]},
            },
        ])).toEqual([
            {
                name: "My albums",
                criterion: {
                    selfOwned: true,
                    owners: [],
                },
                avatars: [currentUser.picture],
            },
            {
                name: "All albums",
                criterion: {
                    owners: []
                },
                avatars: [currentUser.picture],
            },
            {
                name: herselfUser.name,
                criterion: {
                    owners: [herselfOwner]
                },
                avatars: [],
            },
        ])
    })

    it("should return the options 'All albums', 'Her albums', 'His albums' if all albums are owned by two other owners ; pictures are ordered by owner name", () => {
        expect(generateAlbumFilterOptions(currentUser, [
            {
                ...janBaseAlbum,
                albumId: {owner: himselfOwner, folderName: "jan-25"},
                ownedBy: himself,
            },
            {
                ...febBaseAlbum,
                albumId: {owner: herselfOwner, folderName: "jan-25"},
                ownedBy: {name: herselfUser.name, users: [herselfUser]},
            },
        ])).toEqual([
            {
                name: "All albums",
                criterion: {
                    owners: []
                },
                avatars: [currentUser.picture, herselfUser.picture, "a-his-face.jpg"],
            },
            {
                name: herselfUser.name,
                criterion: {
                    owners: [herselfOwner]
                },
                avatars: [herselfUser.picture],
            },
            {
                name: himself.name,
                criterion: {
                    owners: [himselfOwner]
                },
                avatars: [himself.users[0].picture],
            },
        ])
    })

    it("should return the options 'All albums', 'Her albums', 'His albums' in alphabetic order when albums comes in a different order", () => {
        expect(generateAlbumFilterOptions(currentUser, [
            {
                ...febBaseAlbum,
                albumId: {owner: herselfOwner, folderName: "jan-25"},
                ownedBy: {name: herselfUser.name, users: [herselfUser]},
            },
            {
                ...janBaseAlbum,
                albumId: {owner: himselfOwner, folderName: "jan-25"},
                ownedBy: himself,
            },
        ])).toEqual([
            {
                name: "All albums",
                criterion: {
                    owners: []
                },
                avatars: [currentUser.picture, herselfUser.picture, himself.users[0].picture],
            },
            {
                name: herselfUser.name,
                criterion: {
                    owners: [herselfOwner]
                },
                avatars: [herselfUser.picture],
            },
            {
                name: himself.name,
                criterion: {
                    owners: [himselfOwner]
                },
                avatars: [himself.users[0].picture],
            },
        ])
    })
})

export {}