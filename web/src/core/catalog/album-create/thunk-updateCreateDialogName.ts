import {createSimpleThunkDeclaration} from "src/libs/dthunks";
import {createDialogNameChanged} from "./action-createDialogNameChanged";

export const updateCreateDialogNameDeclaration = createSimpleThunkDeclaration(createDialogNameChanged);
