import {CatalogViewerState} from "./language";
import {Action, createGenericReducer} from "@/libs/daction";

export * from "./album-create/selector-createDialogSelector";
export * from "./album-delete/selector-deleteDialogSelector";
export * from "./album-edit-dates/selector-editDatesDialogSelector";
export * from "./album-edit-name/selector-editNameDialogSelector";
export * from "./sharing/selector-sharingDialogSelector";

export type CatalogViewerAction = Action<CatalogViewerState, any>


function createCatalogReducer(): (state: CatalogViewerState, action: Action<CatalogViewerState> | CatalogViewerAction) => CatalogViewerState {
    return createGenericReducer();
}

export const catalogReducer = createCatalogReducer();
export {editNameDialogNoError} from "./language";