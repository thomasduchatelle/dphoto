import {CatalogViewerState} from "../language";
import {editDatesDialogClosed} from "./action-editDatesDialogClosed";
import {ThunkDeclaration} from "../../thunk-engine";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {ActionWithReducer} from "../common/action-factory";

export function closeEditDatesDialogThunk(
    dispatch: (action: ActionWithReducer) => void
): void {
    dispatch(editDatesDialogClosed());
}

export const closeEditDatesDialogDeclaration: ThunkDeclaration<
    CatalogViewerState,
    {},
    () => void,
    CatalogFactoryArgs
> = {
    selector: () => ({}),
    factory: ({dispatch}) => {
        return closeEditDatesDialogThunk.bind(null, dispatch);
    },
};
