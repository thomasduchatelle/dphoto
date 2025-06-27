import {createSimpleThunkDeclaration} from "src/libs/dthunks";
import {editNameDialogOpened} from "./action-editNameDialogOpened";

export const openEditNameDialogDeclaration = createSimpleThunkDeclaration(editNameDialogOpened);
