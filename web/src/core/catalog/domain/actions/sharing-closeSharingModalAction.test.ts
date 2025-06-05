import {closeSharingModalAction, reduceCloseSharingModal} from "./sharing-closeSharingModalAction";
import {herselfUser, loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {sharingDialogSelector} from "./selector-sharingDialogSelector";

describe("reduceCloseSharingModal", () => {
    it("should close the sharing modal by clearing the shareModel property", () => {
        const initial = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                    }
                ],
            }
        };
        const action = closeSharingModalAction();
        expect(sharingDialogSelector(reduceCloseSharingModal(initial, action))).toEqual({
            open: false,
            sharedWith: [],
        });
    });
});
