import {createSimpleThunkDeclaration} from "@/libs/dthunks";
import {customFolderNameToggled} from "./action-customFolderNameToggled";

export const changeFolderNameEnabledDeclaration = createSimpleThunkDeclaration(customFolderNameToggled);
