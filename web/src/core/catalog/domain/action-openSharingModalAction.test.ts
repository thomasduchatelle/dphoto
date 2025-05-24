import {reduceOpenSharingModal, openSharingModalAction, OpenSharingModalAction} from "./action-openSharingModalAction";
import {SharingType, UserDetails, CatalogViewerState} from "./catalog-state";

describe("reduceOpenSharingModal", () => {
    const herselfUser: UserDetails = {email: "her@self.com", name: "Herself", picture: "her-face.jpg"};
    const herselfOwner = "herself";

    const twoAlbums = [
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

    const loadedStateWithTwoAlbums: CatalogViewerState = {
        currentUser: {picture: "my-face.jpg"},
        allAlbums: twoAlbums,
        albumFilterOptions: [],
        albumFilter: {criterion: {owners: []}, avatars: [], name: ""},
        albums: twoAlbums,
        medias: [],
        albumNotFound: false,
        albumsLoaded: true,
        mediasLoaded: true,
    };

    it("should open the sharing modal with the appropriate albumId and already-shared list", () => {
        const action = openSharingModalAction({
            albumId: twoAlbums[0].albumId,
        });

        const expected = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                        role: SharingType.visitor,
                    }
                ],
            }
        };
        expect(reduceOpenSharingModal(loadedStateWithTwoAlbums, action)).toEqual(expected);
    });
    it("should close the sharing modal by clearing the shareModel property if album not found", () => {
        const action = openSharingModalAction({
            albumId: {owner: "notfound", folderName: "notfound"},
        });

        const expected = {
            ...loadedStateWithTwoAlbums,
            shareModal: undefined,
        };
        expect(reduceOpenSharingModal(loadedStateWithTwoAlbums, action)).toEqual(expected);
    });
});
