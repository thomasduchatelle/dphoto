import {closeSharingModalAction, reduceCloseSharingModal} from "./action-closeSharingModalAction";
import {loadedStateWithTwoAlbums, twoAlbums, herselfUser} from "./tests/test-helper-state";

describe("reduceCloseSharingModal", () => {
    it("should close the sharing modal by clearing the shareModel property", () => {
        const initial = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                        role: "visitor",
                    }
                ],
            }
        };
        const action = closeSharingModalAction();
        const expected = loadedStateWithTwoAlbums;
        expect(reduceCloseSharingModal(initial, action)).toEqual(expected);
    });
});
