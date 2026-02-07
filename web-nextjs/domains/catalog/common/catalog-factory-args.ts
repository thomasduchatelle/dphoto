import {Action} from "@/libs/daction";
import {CatalogViewerState} from "../language";

export interface CatalogFactoryArgs<Adapter> {
    adapter: Adapter;
    dispatch: (action: Action<CatalogViewerState, any>) => void;
}
