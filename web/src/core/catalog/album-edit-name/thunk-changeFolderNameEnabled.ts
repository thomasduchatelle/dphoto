import {createSimpleThunkDeclaration} from "src/libs/dthunks";
import {folderNameEnabledChanged} from "./action-folderNameEnabledChanged";

export const changeFolderNameEnabledDeclaration = createSimpleThunkDeclaration(folderNameEnabledChanged);
