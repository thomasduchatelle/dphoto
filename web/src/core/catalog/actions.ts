import {CatalogViewerState} from "./language";
import {Action, createGenericReducer} from "src/libs/daction";

export * from "./album-create/selector-createDialogSelector";
export * from "./album-delete/selector-deleteDialogSelector";
export * from "./album-edit-dates/selector-editDatesDialogSelector";
export * from "./sharing/selector-sharingDialogSelector";

export type CatalogViewerAction = Action<CatalogViewerState, any>


function createCatalogReducer(): (state: CatalogViewerState, action: Action<CatalogViewerState> | CatalogViewerAction) => CatalogViewerState {
    return createGenericReducer();
}

export const catalogReducer = createCatalogReducer();
