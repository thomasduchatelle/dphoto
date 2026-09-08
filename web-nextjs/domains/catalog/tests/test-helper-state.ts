import {
    Album,
    AlbumFilterEntry,
    CatalogViewerState,
    CreateDialog,
    CurrentUserInsight,
    DeleteDialog,
    EditDatesDialog,
    EditNameDialog,
    editNameDialogNoError,
    MediaType,
    MediaWithinADay,
    UserDetails
} from "../language";
import {AlbumListActionsProps, CatalogViewerPageSelection} from "../navigation";
import {EditNameDialogSelection} from "../album-edit-name";
import {CreateDialogSelection} from "../album-create";

// **IMPORTANT** - to LLM Agents
// Use the constants defined in this file in all your tests to make them more readable, and robust to changes
// Update this file **only if you add a new property** to set a sensible default value

// use myselfUser as a default and current user
export const myselfUser: CurrentUserInsight = {picture: "my-face.jpg", isOwner: true};

// use herselfUser when a second user is required
export const herselfUser: UserDetails = {email: "her@self.com", name: "Herself", picture: "her-face.jpg"};
export const herselfOwner = "herself";

// use twoAlbums as default albums
// - index `0` is January 2025 directly owned and shared to 'herself'
// - index `1` is February 2025 owned by herself, and shared to the current user 'myself'
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

// use march2025 as a third album when required - it is not loaded by default
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

// use sampleAlbums when realistic albums with thumbnails are required (storybook stories) - index 0 is the newest
export const sampleAlbums: Album[] = [
    {
        albumId: {owner: "sandfall", folderName: "clair-obscur"},
        name: "Clair Obscur",
        start: new Date(2025, 3, 24),
        end: new Date(2025, 5, 1),
        totalCount: 47,
        temperature: 6.7,
        relativeTemperature: 1,
        sharedWith: [],
        thumbnails: [
            "/thumbnails/clair-obscur-1.jpg",
            "/thumbnails/clair-obscur-2.jpg",
            "/thumbnails/clair-obscur-3.jpg",
            "/thumbnails/clair-obscur-4.jpg",
        ],
    },
    {
        albumId: {owner: "sony", folderName: "astro-bot"},
        name: "Astro Bot",
        start: new Date(2024, 8, 6),
        end: new Date(2024, 8, 12),
        totalCount: 23,
        temperature: 12.1,
        relativeTemperature: 0.72,
        sharedWith: [{user: {name: "Tony Stark", email: "ironman@avenger.com", picture: "/tonystark-profile.jpg"}}],
        thumbnails: [
            "/thumbnails/astro-bot-01.jpg",
            "/thumbnails/astro-bot-02.jpg",
            "/thumbnails/astro-bot-03.jpg",
        ],
    },
    {
        albumId: {owner: "cdprojekt", folderName: "the-witcher-3"},
        name: "The Witcher 3: Wild Hunt",
        start: new Date(2024, 4, 30),
        end: new Date(2024, 5, 14),
        totalCount: 189,
        temperature: 9.4,
        relativeTemperature: 0.55,
        sharedWith: [],
        thumbnails: [
            "/thumbnails/the-witcher-3-01.jpg",
            "/thumbnails/the-witcher-3-02.jpg",
            "/thumbnails/the-witcher-3-03.jpg",
        ],
    },
    {
        albumId: {owner: "kojima", folderName: "death-stranding-1"},
        name: "Death Stranding",
        start: new Date(2023, 10, 8),
        end: new Date(2023, 10, 22),
        totalCount: 312,
        temperature: 17.5,
        relativeTemperature: 1.0,
        ownedBy: {name: "Kojima", users: [{name: "Tony Stark", email: "ironman@avenger.com", picture: "/tonystark-profile.jpg"}]},
        sharedWith: [],
        thumbnails: [
            "/thumbnails/death-stranding-1-01.jpg",
            "/thumbnails/death-stranding-1-02.jpg",
            "/thumbnails/death-stranding-1-03.jpg",
            "/thumbnails/death-stranding-1-04.jpg",
        ],
    },
    {
        albumId: {owner: "kojima", folderName: "death-stranding-2"},
        name: "Death Stranding 2: On The Beach",
        start: new Date(2025, 5, 5),
        end: new Date(2025, 5, 30),
        totalCount: 78,
        temperature: 4.1,
        relativeTemperature: 0.22,
        sharedWith: [{user: {name: "Tony Stark", email: "ironman@avenger.com", picture: "/tonystark-profile.jpg"}}],
        thumbnails: ["/thumbnails/death-stranding-2-01.jpg"],
    },
    {
        albumId: {owner: "sandfall", folderName: "clair-obscur-dlc"},
        name: "Clair Obscur DLC",
        start: new Date(2025, 8, 10),
        end: new Date(2025, 8, 15),
        totalCount: 11,
        temperature: 1.2,
        relativeTemperature: 0.06,
        sharedWith: [],
        thumbnails: [],
    },
]

export const twoAlbumsNoFilterOptions: AlbumFilterEntry = {
    criterion: {
        owners: []
    },
    avatars: [`${myselfUser.picture}`, `${herselfUser.picture}`],
    name: "All albums",
};

// use someMedias as default medias returned from the adapters
const someMedias = [{
    id: "media-1",
    type: MediaType.IMAGE,
    time: new Date("2025-01-05T12:42:00Z"),
    uiRelativePath: "media-1/image.jpg",
    contentPath: "/media-1.jpg",
    thumbnailUrl: "/media-1.jpg?w=360",
    source: "",
}];

// use someMediasByDays as what is expected in the state when `someMedias` are received by the adapter
export const someMediasByDays: MediaWithinADay[] = [{
    day: new Date(2025, 0, 5),
    medias: someMedias,
}]

// use it as a default ready state: it always reflects how a loaded page looks like
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
    medias: someMediasByDays,
    albumNotFound: false,
    mediasLoadedFromAlbumId: twoAlbums[0].albumId,
    albumsLoaded: true,
    mediasLoaded: true,
};

// use it as default selection - it matches the loaded state through the selectors
export const selectionForLoadedStateWithTwoAlbums: CatalogViewerPageSelection = {
    albumsLoaded: true,
    albums: twoAlbums,
    displayedAlbum: twoAlbums[0],
    previousAlbum: twoAlbums[1],
    nextAlbum: undefined,
    medias: someMediasByDays,
    mediasLoaded: true,
    albumNotFound: false,
};

// use it as a default opened delete dialog - it match what would be expected from the `loadedStateWithTwoAlbums`
export const deleteDialogWithOneAlbum: DeleteDialog = {
    type: "DeleteDialog",
    deletableAlbums: [twoAlbums[0]],
    initialSelectedAlbumId: twoAlbums[0].albumId,
    isLoading: false,
};

// use it as default opened create dialog - it match the March 2025 album
export const createDialogPrefilledForMar25: CreateDialog = {
    type: "CreateDialog",
    albumName: "March 2025",
    customFolderName: "",
    isCustomFolderNameEnabled: false,
    nameError: {},
    startDate: new Date("2025-02-01Z"),
    endDate: new Date("2025-03-31Z"),
    startAtDayStart: true,
    endAtDayEnd: true,
    isLoading: false,
}

export const createDialogSelectionPrefilledForMar25: CreateDialogSelection = {
    albumName: "March 2025",
    canSubmit: true,
    end: new Date("2025-03-31Z"),
    endsAtEndOfTheDay: true,
    customFolderName: "",
    isCustomFolderNameEnabled: false,
    isLoading: false,
    open: true,
    start: new Date("2025-02-01Z"),
    startsAtStartOfTheDay: true,

}

// use it as default opened edit dates dialog - it matches the January 2025 album
export const editDatesDialogForJanAlbum: EditDatesDialog = {
    type: "EditDatesDialog",
    albumId: twoAlbums[0].albumId,
    albumName: twoAlbums[0].name,
    startDate: twoAlbums[0].start,
    endDate: twoAlbums[0].end,
    startAtDayStart: true,
    endAtDayEnd: true,
    isLoading: false,
}

export const editJanAlbumNameDialog: EditNameDialog = {
    type: "EditNameDialog",
    albumId: twoAlbums[0].albumId,
    albumName: twoAlbums[0].name,
    originalFolderName: twoAlbums[0].albumId.folderName,
    customFolderName: "",
    isCustomFolderNameEnabled: false,
    isLoading: false,
    nameError: editNameDialogNoError,
}

export const editJanAlbumNameSelection: EditNameDialogSelection = {
    isOpen: true,
    albumName: twoAlbums[0].name,
    originalName: twoAlbums[0].name,
    customFolderName: "",
    isCustomFolderNameEnabled: false,
    isLoading: false,
    isSaveEnabled: true,
}

export const albumListActionsPropsForLoadedState: AlbumListActionsProps = {
    albumFilter: twoAlbumsNoFilterOptions,
    albumFilterOptions: loadedStateWithTwoAlbums.albumFilterOptions,
    displayedAlbumIdIsOwned: true,
    hasAlbumsToDelete: true,
    canCreateAlbums: true,
}




