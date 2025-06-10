import {closeDeleteAlbumDialogThunk} from "./thunk-closeDeleteAlbumDialog";
import {deleteAlbumDialogClosed} from "./action-deleteAlbumDialogClosed";

describe("closeDeleteAlbumDialogThunk", () => {
    it("should dispatch closeDeleteAlbumDialogAction", () => {
        const dispatched: any[] = [];
        closeDeleteAlbumDialogThunk(dispatched.push.bind(dispatched));
        expect(dispatched).toEqual([
            deleteAlbumDialogClosed()
        ]);
    });
});
