import {openSharingModalAction, reduceOpenSharingModal} from "./sharing-openSharingModalAction";
import {herselfUser, loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {SharingDialogFrag, sharingDialogSelector} from "./selector-sharingDialogSelector";

describe("reduceOpenSharingModal", () => {

    it("should open the sharing modal with the appropriate albumId and already-shared list", () => {
        const action = openSharingModalAction(twoAlbums[0].albumId);

        const expected: SharingDialogFrag = {
            open: true,
            sharedWith: [
                {
                    user: herselfUser,
                }
            ],
        };
        expect(sharingDialogSelector(reduceOpenSharingModal(loadedStateWithTwoAlbums, action))).toEqual(expected);
    });

    it("should close the sharing modal by clearing the shareModel property if album not found", () => {
        const action = openSharingModalAction({
            albumId: {owner: "notfound", folderName: "notfound"},
        });

        const expected: SharingDialogFrag = {
            open: false,
            sharedWith: [],
        };
        expect(sharingDialogSelector(reduceOpenSharingModal(loadedStateWithTwoAlbums, action))).toEqual(expected);
    });
});
