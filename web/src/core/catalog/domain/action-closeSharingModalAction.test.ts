import {closeSharingModalAction, reduceCloseSharingModal} from "./action-closeSharingModalAction";
import {herselfUser, loadedStateWithTwoAlbums, twoAlbums} from "./tests/test-helper-state";
import {SharingType} from "./catalog-state";

describe("reduceCloseSharingModal", () => {
    it("should close the sharing modal by clearing the shareModel property", () => {
        const initial = {
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
        const action = closeSharingModalAction();
        expect(reduceCloseSharingModal(initial, action)).toEqual(loadedStateWithTwoAlbums);
    });
});
