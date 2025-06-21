import {createSimpleThunkDeclaration} from "../../thunk-engine/simple-thunk-factory";
import {editDatesDialogOpened} from "./action-editDatesDialogOpened";

export const openEditDatesDialogDeclaration = createSimpleThunkDeclaration(editDatesDialogOpened);
