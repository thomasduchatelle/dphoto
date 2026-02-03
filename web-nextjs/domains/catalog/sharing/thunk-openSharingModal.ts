import {sharingModalOpened} from "./action-sharingModalOpened";
import {createSimpleThunkDeclaration} from "@/libs/dthunks";

export const openSharingModalDeclaration = createSimpleThunkDeclaration(sharingModalOpened);