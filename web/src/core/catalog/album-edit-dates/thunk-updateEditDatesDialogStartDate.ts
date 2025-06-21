import {createSimpleThunkDeclaration} from "../../thunk-engine/simple-thunk-factory";
import {editDatesDialogStartDateUpdated} from "./action-editDatesDialogStartDateUpdated";

export const updateEditDatesDialogStartDateDeclaration = createSimpleThunkDeclaration(editDatesDialogStartDateUpdated);
