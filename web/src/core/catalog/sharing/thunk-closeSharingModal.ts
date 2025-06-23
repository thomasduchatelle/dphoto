import {CatalogViewerState} from "../language";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {SharingModalClosed, sharingModalClosed} from "./action-sharingModalClosed";
import {createSimpleThunkDeclaration, ThunkDeclaration} from "src/libs/dthunks";

export function closeSharingModalThunk(dispatch: (action: SharingModalClosed) => void): void {
    dispatch(sharingModalClosed());
}

export const closeSharingModalDeclaration = createSimpleThunkDeclaration(sharingModalClosed);
