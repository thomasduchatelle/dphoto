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
import {CatalogViewerPageSelection} from "../navigation";
import {EditNameDialogSelection} from "../album-edit-name";

// **IMPORTANT** - to LLM Agents
// Use the constants defined in this file in all your tests to make them more readable, and robust to changes
// Update this file **only if you add a new property** to set a sensible default value

// use myselfUser as a default and current user
export const myselfUser: CurrentUserInsight = {picture: "my-face.jpg"};

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
    albumFilter: loadedStateWithTwoAlbums.albumFilter,
    albumFilterOptions: loadedStateWithTwoAlbums.albumFilterOptions,
    albumsLoaded: true,
    albums: twoAlbums,
    displayedAlbum: twoAlbums[0],
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
    name: "March 2025",
    startDate: new Date("2025-02-01"),
    endDate: new Date("2025-03-31"),
    startAtDayStart: true,
    endAtDayEnd: true,
    forceFolderName: "",
    withCustomFolderName: false,
    isLoading: false,
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
    customFolderName: "",
    isCustomFolderNameEnabled: false,
    isLoading: false,
    error: editNameDialogNoError,
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




