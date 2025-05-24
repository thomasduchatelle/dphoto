import {albumMatchCriterion, CatalogViewerState, CurrentUserInsight, Sharing} from "./catalog-state";
import {CatalogViewerAction} from "./catalog-actions";
import {albumIdEquals} from "./utils-albumIdEquals";
import {albumFilterAreCriterionEqual, ALL_ALBUMS_FILTER_CRITERION, DEFAULT_ALBUM_FILTER_ENTRY, generateAlbumFilterOptions} from "./catalog-common-modifiers";

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

export const catalogReducerFunction = (current: CatalogViewerState, action: CatalogViewerAction): CatalogViewerState => {
    switch (action.type) {
        case "StartLoadingMediasAction":
            return {
                ...current,
                medias: [],
                error: undefined,
                loadingMediasFor: action.albumId,
                albumNotFound: false,
                mediasLoaded: false,
            }

        case "MediasLoadedAction":
            if (current.loadingMediasFor && !albumIdEquals(current.loadingMediasFor, action.albumId)) {
                // concurrency management - ignore if not the last album requested
                return current
            }

            return {
                ...current,
                loadingMediasFor: undefined,
                mediasLoadedFromAlbumId: action.albumId,
                medias: action.medias,
                error: undefined,
                mediasLoaded: true,
                albumNotFound: false,
            }

        case "NoAlbumAvailableAction":
            return {
                currentUser: current.currentUser,
                albumNotFound: true,
                allAlbums: [],
                albums: [],
                medias: [],
                albumsLoaded: true,
                mediasLoaded: true,
                albumFilterOptions: [DEFAULT_ALBUM_FILTER_ENTRY],
                albumFilter: DEFAULT_ALBUM_FILTER_ENTRY,
            }

        case "AlbumsAndMediasLoadedAction":
            const albumFilterOptions = generateAlbumFilterOptions(current.currentUser, action.albums);

            return {
                currentUser: current.currentUser,
                albumNotFound: false,
                allAlbums: action.albums,
                albums: action.albums,
                mediasLoadedFromAlbumId: action.selectedAlbum?.albumId,
                medias: action.medias,
                albumsLoaded: true,
                mediasLoaded: true,
                albumFilterOptions,
                albumFilter: albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, ALL_ALBUMS_FILTER_CRITERION)) ?? DEFAULT_ALBUM_FILTER_ENTRY
            }

        case "AlbumsFilteredAction":
            const filteredAlbums = current.allAlbums.filter(albumMatchCriterion(action.criterion))

            const allAlbumFilter = current.albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, ALL_ALBUMS_FILTER_CRITERION)) ?? DEFAULT_ALBUM_FILTER_ENTRY

            return {
                ...current,
                albums: filteredAlbums,
                albumFilter: current.albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, action.criterion)) ?? allAlbumFilter,
            }

        case "MediaFailedToLoadAction": {
            const allAlbums = action.albums ?? current.allAlbums
            const albumFilterOptions = generateAlbumFilterOptions(current.currentUser, allAlbums)
            const albumFilter = albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, current.albumFilter.criterion)) ?? DEFAULT_ALBUM_FILTER_ENTRY
            const albums = allAlbums.filter(albumMatchCriterion(albumFilter.criterion))

            return {
                currentUser: current.currentUser,
                allAlbums,
                albumFilterOptions,
                albumFilter,
                albums,
                albumNotFound: false,
                medias: [],
                error: action.error,
                albumsLoaded: true,
                mediasLoaded: true,
            }
        }

        case "AlbumsLoadedAction": {
            const albumFilterOptions = generateAlbumFilterOptions(current.currentUser, action.albums)
            const albumFilter = albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, current.albumFilter.criterion)) ?? DEFAULT_ALBUM_FILTER_ENTRY

            let staging: CatalogViewerState = {
                ...current,
                albumFilterOptions,
                albumFilter,
                allAlbums: action.albums,
                albums: action.albums.filter(albumMatchCriterion(current.albumFilter.criterion)),
                error: undefined,
                albumsLoaded: true,
            }

            if (!staging.albums.find(album => albumIdEquals(album.albumId, action.redirectTo))) {
                const albumFilter = current.albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, DEFAULT_ALBUM_FILTER_ENTRY.criterion)) ?? DEFAULT_ALBUM_FILTER_ENTRY
                staging = {
                    ...staging,
                    albumFilter,
                    albums: action.albums.filter(albumMatchCriterion(albumFilter.criterion)),
                }
            }

            return staging
        }

        case "OpenSharingModalAction": {
            // Find the album in the current state
            const album = current.allAlbums.find(a => albumIdEquals(a.albumId, action.albumId));
            return {
                ...current,
                shareModal: album
                    ? {
                        sharedAlbumId: album.albumId,
                        sharedWith: sortSharings([...album.sharedWith]),
                    }
                    : undefined,
            };
        }

        case "CloseSharingModalAction": {
            // Remove the shareModal property
            const {shareModal, ...rest} = current;
            return rest as CatalogViewerState;
        }

        case "AddSharingAction": {
            if (!current.shareModal) return current;
            // Replace if user already exists (by email), else add
            const newSharing = action.sharing;
            const updatedSharedWith = [
                ...current.shareModal.sharedWith.filter(s => s.user.email !== newSharing.user.email),
                newSharing
            ];
            return {
                ...current,
                shareModal: {
                    ...current.shareModal,
                    sharedWith: sortSharings(updatedSharedWith),
                }
            };
        }

        case "RemoveSharingAction": {
            if (!current.shareModal) return current;
            const updatedSharedWith = current.shareModal.sharedWith.filter(
                s => s.user.email !== action.email
            );
            return {
                ...current,
                shareModal: {
                    ...current.shareModal,
                    sharedWith: updatedSharedWith,
                }
            };
        }

        case "SharingModalErrorAction": {
            if (!current.shareModal) return current;
            return {
                ...current,
                shareModal: {
                    ...current.shareModal,
                    error: action.error,
                }
            };
        }

        default:
            return current
    }
}

function sortSharings(sharings: Sharing[]): Sharing[] {
    return sharings.slice().sort((a, b) => {
        const nameA = a.user.name?.trim() || "";
        const nameB = b.user.name?.trim() || "";
        if (nameA && nameB) {
            const cmp = nameA.localeCompare(nameB);
            if (cmp !== 0) return cmp;
            return a.user.email.localeCompare(b.user.email);
        }
        if (!nameA && !nameB) {
            return a.user.email.localeCompare(b.user.email);
        }
        if (!nameA) return 1;
        if (!nameB) return -1;
        return 0;
    });
}
