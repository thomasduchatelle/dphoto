import {closeSharingModalThunk} from "./thunk-closeSharingModal";
import {SharingModalClosed, sharingModalClosed} from "./action-sharingModalClosed";

describe("onCloseThunk", () => {
    it("should dispatches closeSharingModalAction when called", () => {
        const dispatched: SharingModalClosed[] = [];
        closeSharingModalThunk(dispatched.push.bind(dispatched));
        expect(dispatched).toEqual([
            sharingModalClosed()
        ]);
    });
});
