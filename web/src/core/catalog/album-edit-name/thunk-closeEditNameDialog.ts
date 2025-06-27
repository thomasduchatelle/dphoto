import {createSimpleThunkDeclaration} from "src/libs/dthunks";
import {editNameDialogClosed} from "./action-editNameDialogClosed";

export const closeEditNameDialogDeclaration = createSimpleThunkDeclaration(editNameDialogClosed);
