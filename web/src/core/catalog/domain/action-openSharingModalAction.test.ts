import {openSharingModalAction, reduceOpenSharingModal} from "./action-openSharingModalAction";
import {SharingType} from "./catalog-state";
import {herselfUser, loadedStateWithTwoAlbums, twoAlbums} from "./tests/test-helper-state";

describe("reduceOpenSharingModal", () => {

    it("should open the sharing modal with the appropriate albumId and already-shared list", () => {
        const action = openSharingModalAction(twoAlbums[0].albumId);

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
