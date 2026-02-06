import {createSimpleThunkDeclaration} from "@/libs/dthunks";
import {createDialogClosed} from "./action-createDialogClosed";

export const closeCreateDialogDeclaration = createSimpleThunkDeclaration(createDialogClosed);
