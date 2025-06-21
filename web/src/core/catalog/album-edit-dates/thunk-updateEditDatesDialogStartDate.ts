import {ThunkDeclaration} from "../../thunk-engine";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {CatalogViewerState} from "../language";
import {editDatesDialogStartDateUpdated, EditDatesDialogStartDateUpdated} from "./action-editDatesDialogStartDateUpdated";

export async function updateEditDatesDialogStartDateThunk(
    dispatch: (action: EditDatesDialogStartDateUpdated) => void,
    startDate: Date | null
): Promise<void> {
    if (startDate) {
        dispatch(editDatesDialogStartDateUpdated(startDate));
    }
}

export const updateEditDatesDialogStartDateDeclaration: ThunkDeclaration<
    CatalogViewerState,
    {},
    (startDate: Date | null) => Promise<void>,
    CatalogFactoryArgs
> = {
    selector: (state: CatalogViewerState) => ({}),
    factory: ({dispatch}) => {
        return (startDate: Date | null) => updateEditDatesDialogStartDateThunk(dispatch, startDate);
    },
};
