import {editDatesDialogClosed} from "./action-editDatesDialogClosed";
import {createSimpleThunkDeclaration} from "@/libs/dthunks";

export const closeEditDatesDialogDeclaration = createSimpleThunkDeclaration(editDatesDialogClosed);
