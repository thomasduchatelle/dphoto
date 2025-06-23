import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {CatalogViewerState} from "../language";
import {editDatesDialogEndDateUpdated} from "./action-editDatesDialogEndDateUpdated";
import {ThunkDeclaration, createSimpleThunkDeclaration} from "src/libs/dthunks";

export const updateEditDatesDialogEndDateDeclaration: ThunkDeclaration<
    CatalogViewerState,
    {},
    (endDate: Date | null) => Promise<void>,
    CatalogFactoryArgs
> = createSimpleThunkDeclaration(editDatesDialogEndDateUpdated);
