import {createSimpleThunkDeclaration} from "../../thunk-engine/simple-thunk-factory";
import {editDatesDialogClosed} from "./action-editDatesDialogClosed";

export const closeEditDatesDialogDeclaration = createSimpleThunkDeclaration(editDatesDialogClosed);
