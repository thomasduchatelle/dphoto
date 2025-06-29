import {createSimpleThunkDeclaration} from "src/libs/dthunks";
import {customFolderNameToggled} from "./action-customFolderNameToggled";

export const changeFolderNameEnabledDeclaration = createSimpleThunkDeclaration(customFolderNameToggled);
