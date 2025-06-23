import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {CatalogViewerState} from "../language";
import {editDatesDialogStartDateUpdated} from "./action-editDatesDialogStartDateUpdated";
import {ThunkDeclaration, createSimpleThunkDeclaration} from "src/libs/dthunks";

export const updateEditDatesDialogStartDateDeclaration: ThunkDeclaration<
    CatalogViewerState,
    {},
    (startDate: Date | null) => Promise<void>,
    CatalogFactoryArgs
> = createSimpleThunkDeclaration(editDatesDialogStartDateUpdated);
