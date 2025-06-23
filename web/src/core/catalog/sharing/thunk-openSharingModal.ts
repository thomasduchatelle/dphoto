import {sharingModalOpened} from "./action-sharingModalOpened";
import {createSimpleThunkDeclaration} from "src/libs/dthunks";

export const openSharingModalDeclaration = createSimpleThunkDeclaration(sharingModalOpened);