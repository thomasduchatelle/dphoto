import {createSimpleThunkDeclaration} from "src/libs/dthunks";
import {createDialogClosed} from "./action-createDialogClosed";

export const closeCreateDialogDeclaration = createSimpleThunkDeclaration(createDialogClosed);
