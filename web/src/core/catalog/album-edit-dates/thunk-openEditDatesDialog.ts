import {CatalogViewerState} from "../language";
import {editDatesDialogOpened, EditDatesDialogOpened} from "./action-editDatesDialogOpened";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {ThunkDeclaration, createSimpleThunkDeclaration} from "src/libs/dthunks";

export const openEditDatesDialogDeclaration: ThunkDeclaration<
    CatalogViewerState,
    {},
    () => void,
    CatalogFactoryArgs
> = createSimpleThunkDeclaration(editDatesDialogOpened);
