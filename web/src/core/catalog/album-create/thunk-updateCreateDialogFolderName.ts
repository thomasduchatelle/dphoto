import {createSimpleThunkDeclaration} from "src/libs/dthunks";
import {createDialogFolderNameChanged} from "./action-createDialogFolderNameChanged";

export const updateCreateDialogFolderNameDeclaration = createSimpleThunkDeclaration(createDialogFolderNameChanged);
