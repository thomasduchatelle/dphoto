import {editDatesDialogOpened} from "./action-editDatesDialogOpened";
import {createSimpleThunkDeclaration} from "@/libs/dthunks";

export const openEditDatesDialogDeclaration = createSimpleThunkDeclaration(editDatesDialogOpened);
