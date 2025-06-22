import {CatalogViewerState} from "../language";
import {editDatesDialogClosed, EditDatesDialogClosed} from "./action-editDatesDialogClosed";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {ThunkDeclaration} from "src/libs/dthunks";

export function closeEditDatesDialogThunk(
    dispatch: (action: EditDatesDialogClosed) => void
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
