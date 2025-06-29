import {createSimpleThunkDeclaration} from "src/libs/dthunks";
import {customFolderNameChanged} from "./action-customFolderNameChanged";

export const changeFolderNameDeclaration = createSimpleThunkDeclaration(customFolderNameChanged);
