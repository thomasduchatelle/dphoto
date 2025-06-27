import {createSimpleThunkDeclaration} from "src/libs/dthunks";
import {folderNameChanged} from "./action-folderNameChanged";

export const changeFolderNameDeclaration = createSimpleThunkDeclaration(folderNameChanged);
