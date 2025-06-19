import {editAlbumDatesDialogClosed, EditAlbumDatesDialogClosed} from "./action-editAlbumDatesDialogClosed";
import {CatalogViewerState} from "../language";
import {ThunkDeclaration} from "../../thunk-engine";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";

export function closeEditAlbumDatesDialogThunk(
    dispatch: (action: EditAlbumDatesDialogClosed) => void,
): void {
    dispatch(editAlbumDatesDialogClosed());
}

export const closeEditAlbumDatesDialogDeclaration: ThunkDeclaration<
    CatalogViewerState,
    {}, // No specific partial state needed from selector for this thunk
    () => void,
    CatalogFactoryArgs
> = {
    selector: (state: CatalogViewerState) => ({}), // No specific data needed from state for the factory
    factory: ({dispatch, app}) => {
        // Bind the thunk with dispatch
        return () => closeEditAlbumDatesDialogThunk(dispatch);
    },
};
