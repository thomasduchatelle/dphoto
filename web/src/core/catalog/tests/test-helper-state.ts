import {Album, AlbumFilterEntry, CatalogViewerState, CurrentUserInsight, MediaType, MediaWithinADay, UserDetails} from "../language";
import {CatalogViewerPageSelection} from "../navigation";

export const myselfUser: CurrentUserInsight = {picture: "my-face.jpg"};
export const herselfUser: UserDetails = {email: "her@self.com", name: "Herself", picture: "her-face.jpg"};
export const herselfOwner = "herself";

export const twoAlbums: Album[] = [
    {
        albumId: {owner: "myself", folderName: "jan-25"},
        name: "January 2025",
        start: new Date(2025, 0, 1),
        end: new Date(2025, 1, 1),
        totalCount: 42,
        temperature: 0.25,
        relativeTemperature: 1,
        sharedWith: [
            {
                user: herselfUser,
            }
        ],
    },
    {
        albumId: {owner: herselfOwner, folderName: "feb-25"},
        name: "February 2025",
        start: new Date(2025, 1, 1),
        end: new Date(2025, 2, 1),
        totalCount: 12,
        temperature: 0.25,
        relativeTemperature: 1,
        ownedBy: {name: "Herself", users: [herselfUser]},
        sharedWith: [],
    },
]

export const march2025: Album = {
    albumId: {owner: "myself", folderName: "mar-25"},
    name: "March 2025",
    start: new Date(2025, 2, 1),
    end: new Date(2025, 3, 1),
    totalCount: 0,
    temperature: 0,
    relativeTemperature: 0,
    sharedWith: [],
}

export const twoAlbumsNoFilterOptions: AlbumFilterEntry = {
    criterion: {
        owners: []
    },
    avatars: [`${myselfUser.picture}`, `${herselfUser.picture}`],
    name: "All albums",
};

export const someMedias: MediaWithinADay[] = [{
    day: new Date(2025, 0, 5),
    medias: [{
        id: "media-1",
        type: MediaType.IMAGE,
        time: new Date("2025-01-05T12:42:00Z"),
        uiRelativePath: "media-1/image.jpg",
        contentPath: "/media-1.jpg",
        source: "",
    }],
}]

export const loadedStateWithTwoAlbums: CatalogViewerState = {
    currentUser: myselfUser,
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
    mediasLoadedFromAlbumId: twoAlbums[0].albumId,
    albumsLoaded: true,
    mediasLoaded: true,
};

export const selectionForLoadedStateWithTwoAlbums: CatalogViewerPageSelection = {
    albumFilter: loadedStateWithTwoAlbums.albumFilter,
    albumFilterOptions: loadedStateWithTwoAlbums.albumFilterOptions,
    albumsLoaded: true,
    albums: twoAlbums,
    displayedAlbum: twoAlbums[0],
    medias: someMedias,
    mediasLoaded: true,
    mediasLoadedFromAlbumId: twoAlbums[0].albumId,
    loadingMediasFor: undefined,
    albumNotFound: false,
};
