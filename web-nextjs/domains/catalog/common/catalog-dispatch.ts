import {Action} from "@/libs/daction";
import {CatalogViewerState} from "../language";

export interface CatalogDispatch {
    dispatch: (action: Action<CatalogViewerState, any>) => void;
}
