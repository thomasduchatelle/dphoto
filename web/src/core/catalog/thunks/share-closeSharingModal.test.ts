import {closeSharingModalThunk} from "./share-closeSharingModal";
import {catalogActions, CatalogViewerAction} from "../domain";

describe("onCloseThunk", () => {
    it("should dispatches closeSharingModalAction when called", () => {
        const dispatched: CatalogViewerAction[] = [];
        closeSharingModalThunk(dispatched.push.bind(dispatched));
        expect(dispatched).toEqual([
            catalogActions.closeSharingModalAction()
        ]);
    });
});
