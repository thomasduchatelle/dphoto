import {SharingModalClosed, sharingModalClosed} from "./action-sharingModalClosed";
import {createSimpleThunkDeclaration} from "src/libs/dthunks";

export function closeSharingModalThunk(dispatch: (action: SharingModalClosed) => void): void {
    dispatch(sharingModalClosed());
}

export const closeSharingModalDeclaration = createSimpleThunkDeclaration(sharingModalClosed);
