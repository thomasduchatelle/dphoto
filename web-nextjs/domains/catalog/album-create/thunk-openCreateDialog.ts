import {createSimpleThunkDeclaration} from "@/libs/dthunks";
import {createDialogOpened} from "./action-createDialogOpened";

export const openCreateDialogDeclaration = createSimpleThunkDeclaration(createDialogOpened);
