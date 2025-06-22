import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {CatalogViewerState} from "../language";
import {editDatesDialogEndDateUpdated, EditDatesDialogEndDateUpdated} from "./action-editDatesDialogEndDateUpdated";
import {ThunkDeclaration} from "src/libs/thunks";

export async function updateEditDatesDialogEndDateThunk(
    dispatch: (action: EditDatesDialogEndDateUpdated) => void,
    endDate: Date | null
): Promise<void> {
    if (endDate) {
        dispatch(editDatesDialogEndDateUpdated(endDate));
    }
}

export const updateEditDatesDialogEndDateDeclaration: ThunkDeclaration<
    CatalogViewerState,
    {},
    (endDate: Date | null) => Promise<void>,
    CatalogFactoryArgs
> = {
    selector: (state: CatalogViewerState) => ({}),
    factory: ({dispatch}) => {
        return (endDate: Date | null) => updateEditDatesDialogEndDateThunk(dispatch, endDate);
    },
};
