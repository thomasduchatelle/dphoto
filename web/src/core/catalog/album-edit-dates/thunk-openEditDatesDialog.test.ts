import {openEditDatesDialogThunk} from "./thunk-openEditDatesDialog";
import {editDatesDialogOpened} from "./action-editDatesDialogOpened";

describe("thunk:openEditDatesDialog", () => {
    it("dispatches editDatesDialogOpened action", () => {
        const dispatched: any[] = [];

        openEditDatesDialogThunk(dispatched.push.bind(dispatched));

        expect(dispatched).toEqual([
            editDatesDialogOpened()
        ]);
    });
});
