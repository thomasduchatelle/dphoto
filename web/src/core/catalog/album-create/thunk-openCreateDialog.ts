import {createSimpleThunkDeclaration} from "src/libs/dthunks";
import {createDialogOpened} from "./action-createDialogOpened";

export const openCreateDialogDeclaration = createSimpleThunkDeclaration(createDialogOpened);
