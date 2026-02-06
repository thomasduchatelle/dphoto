import {Action} from "@/libs/daction";
import {CatalogViewerState} from "../language";

export interface CatalogFactoryArgs {
    dispatch: (action: Action<CatalogViewerState, any>) => void;
}
