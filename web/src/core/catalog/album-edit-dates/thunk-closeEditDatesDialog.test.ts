import {closeEditDatesDialogThunk} from "./thunk-closeEditDatesDialog";
import {editDatesDialogClosed} from "./action-editDatesDialogClosed";

describe("thunk:closeEditDatesDialog", () => {
    it("dispatches editDatesDialogClosed action", () => {
        const dispatched: any[] = [];

        closeEditDatesDialogThunk(dispatched.push.bind(dispatched));

        expect(dispatched).toEqual([
            editDatesDialogClosed()
        ]);
    });
});
