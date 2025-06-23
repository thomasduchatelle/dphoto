import {CatalogViewerState} from "../language";
import {editDatesDialogClosed, EditDatesDialogClosed} from "./action-editDatesDialogClosed";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {ThunkDeclaration, createSimpleThunkDeclaration} from "src/libs/dthunks";

export const closeEditDatesDialogDeclaration: ThunkDeclaration<
    CatalogViewerState,
    {},
    () => void,
    CatalogFactoryArgs
> = createSimpleThunkDeclaration(editDatesDialogClosed);
