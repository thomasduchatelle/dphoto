import {CatalogViewerState} from "../language";
import {editDatesDialogOpened, EditDatesDialogOpened} from "./action-editDatesDialogOpened";
import {ThunkDeclaration} from "../../thunk-engine";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";

export function openEditDatesDialogThunk(
    dispatch: (action: EditDatesDialogOpened) => void
): void {
    dispatch(editDatesDialogOpened());
}

export const openEditDatesDialogDeclaration: ThunkDeclaration<
    CatalogViewerState,
    {},
    () => void,
    CatalogFactoryArgs
> = {
    selector: () => ({}),
    factory: ({dispatch}) => {
        return openEditDatesDialogThunk.bind(null, dispatch);
    },
};
