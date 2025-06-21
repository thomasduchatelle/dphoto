import {createSimpleThunkDeclaration} from "../../thunk-engine/simple-thunk-factory";
import {editDatesDialogEndDateUpdated} from "./action-editDatesDialogEndDateUpdated";

export const updateEditDatesDialogEndDateDeclaration = createSimpleThunkDeclaration(editDatesDialogEndDateUpdated);
