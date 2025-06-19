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
    {},
    () => void,
    CatalogFactoryArgs
> = {
    selector: ({}: CatalogViewerState) => ({}),
    factory: ({dispatch}) => {
        return closeEditAlbumDatesDialogThunk.bind(null, dispatch);
    },
};
