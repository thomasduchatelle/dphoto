import {createSimpleThunkDeclaration} from "@/libs/dthunks";
import {customFolderNameChanged} from "./action-customFolderNameChanged";

export const changeFolderNameDeclaration = createSimpleThunkDeclaration(customFolderNameChanged);
