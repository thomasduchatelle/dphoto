import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {CatalogViewerState} from "../language";
import {editDatesDialogStartDateUpdated, EditDatesDialogStartDateUpdated} from "./action-editDatesDialogStartDateUpdated";
import {ThunkDeclaration, createSimpleThunkDeclaration} from "src/libs/dthunks";

export const updateEditDatesDialogStartDateDeclaration: ThunkDeclaration<
    CatalogViewerState,
    {},
    (startDate: Date | null) => Promise<void>,
    CatalogFactoryArgs
> = {
    selector: (state: CatalogViewerState) => ({}),
    factory: ({dispatch}) => {
        return async (startDate: Date | null) => {
            if (startDate) {
                dispatch(editDatesDialogStartDateUpdated(startDate));
            }
        };
    },
};
