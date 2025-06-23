import {editDatesDialogClosed} from "./action-editDatesDialogClosed";
import {createSimpleThunkDeclaration} from "src/libs/dthunks";

export const closeEditDatesDialogDeclaration = createSimpleThunkDeclaration(editDatesDialogClosed);
