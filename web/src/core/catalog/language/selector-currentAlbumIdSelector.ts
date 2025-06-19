import {CatalogViewerState} from "./catalog-state";

export function currentAlbumIdSelector({
                                           loadingMediasFor,
                                           mediasLoadedFromAlbumId
                                       }: Pick<CatalogViewerState, "loadingMediasFor"> & Pick<CatalogViewerState, "mediasLoadedFromAlbumId">) {
    return loadingMediasFor || mediasLoadedFromAlbumId;
}