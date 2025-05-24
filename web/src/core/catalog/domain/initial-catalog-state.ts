import {CatalogViewerState, CurrentUserInsight} from "./catalog-state";
import {DEFAULT_ALBUM_FILTER_ENTRY} from "./catalog-common-modifiers";

export const initialCatalogState = (currentUser: CurrentUserInsight): CatalogViewerState => ({
    currentUser,
    albumFilterOptions: [DEFAULT_ALBUM_FILTER_ENTRY],
    albumFilter: DEFAULT_ALBUM_FILTER_ENTRY,
    allAlbums: [],
    albumNotFound: false,
    albums: [],
    medias: [],
    albumsLoaded: false,
    mediasLoaded: false
})