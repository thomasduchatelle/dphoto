import {openDeleteAlbumDialogThunk} from "./thunk-openDeleteAlbumDialog";
import {deleteAlbumDialogOpened} from "./action-deleteAlbumDialogOpened";

describe("openDeleteAlbumDialogThunk", () => {
    it("should dispatch openDeleteAlbumDialogAction", async () => {
        const dispatched: any[] = [];
        openDeleteAlbumDialogThunk(dispatched.push.bind(dispatched));
        expect(dispatched).toEqual([
            deleteAlbumDialogOpened()
        ]);
    });
});
