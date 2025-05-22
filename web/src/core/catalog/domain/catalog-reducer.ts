import {
    Album,
    AlbumFilterCriterion,
    AlbumFilterEntry,
    albumIsOwnedByCurrentUser,
    albumMatchCriterion,
    CatalogViewerState,
    OwnerDetails,
    Sharing
} from "./catalog-state";
import {CatalogViewerAction} from "./catalog-actions";
import {albumIdEquals} from "./utils-albumIdEquals";

const ALL_ALBUMS_FILTER_CRITERION: AlbumFilterCriterion = {owners: []}
const SELF_OWNED_ALBUM_FILTER_CRITERION: AlbumFilterCriterion = {selfOwned: true, owners: []}
const DEFAULT_ALBUM_FILTER_ENTRY: AlbumFilterEntry = {
    criterion: ALL_ALBUMS_FILTER_CRITERION,
    avatars: [],
    name: "All albums",
}

export const initialCatalogState: CatalogViewerState = {
    albumFilterOptions: [DEFAULT_ALBUM_FILTER_ENTRY],
    albumFilter: DEFAULT_ALBUM_FILTER_ENTRY,
    allAlbums: [],
    albumNotFound: false,
    albums: [],
    medias: [],
    albumsLoaded: false,
    mediasLoaded: false
}

export interface CurrentUserInsight {
    picture?: string
}

export const catalogReducerFunction = (currentUser: CurrentUserInsight): (current: CatalogViewerState, action: CatalogViewerAction) => CatalogViewerState => {
    return (current: CatalogViewerState, action: CatalogViewerAction): CatalogViewerState => {
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
                const albumFilterOptions = generateAlbumFilterOptions(currentUser, action.albums);

                return {
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
                const albumFilterOptions = generateAlbumFilterOptions(currentUser, allAlbums)
                const albumFilter = albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, current.albumFilter.criterion)) ?? DEFAULT_ALBUM_FILTER_ENTRY
                const albums = allAlbums.filter(albumMatchCriterion(albumFilter.criterion))

                return {
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
                const albumFilterOptions = generateAlbumFilterOptions(currentUser, action.albums)
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
                const { shareModal, ...rest } = current;
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

export function generateAlbumFilterOptions(currentUser: CurrentUserInsight, albums: Album[]) {
    const currentUserPicture: string[] = currentUser.picture ? [currentUser.picture] : []

    let selfOwnedAlbum = false
    const owners: Map<string, OwnerDetails> = new Map()

    albums.forEach(album => {
        if (albumIsOwnedByCurrentUser(album)) {
            selfOwnedAlbum = true
        } else if (album.ownedBy !== undefined) {
            owners.set(album.albumId.owner, album.ownedBy)
        }
    })

    const options: AlbumFilterEntry[] = []
    if (selfOwnedAlbum && owners.size > 0) {
        options.push({
            name: "My albums",
            criterion: SELF_OWNED_ALBUM_FILTER_CRITERION,
            avatars: currentUserPicture,
        })
    }

    const allPictures: string[] = [...owners.values()].sort(ownersByName).flatMap(owner => owner.users.map(user => user.picture).filter(avatar => avatar !== undefined) as string[])

    options.push({
        name: DEFAULT_ALBUM_FILTER_ENTRY.name,
        criterion: ALL_ALBUMS_FILTER_CRITERION,
        avatars: [...currentUserPicture, ...allPictures],
    })

    if (selfOwnedAlbum || owners.size > 1) {
        [...owners.entries()].sort((a, b) => ownersByName(a[1], b[1])).forEach(owner => {
            options.push({
                name: owner[1].name,
                criterion: {
                    owners: [owner[0]],
                },
                avatars: owner[1].users.map(user => user.picture).filter(avatar => avatar !== undefined) as string[],
            })
        })
    }
    return options;
}

const ownersByName = (a: OwnerDetails, b: OwnerDetails) => a.name.localeCompare(b.name);

function albumFilterAreCriterionEqual(a: AlbumFilterCriterion, b: AlbumFilterCriterion) {
    return a.selfOwned === b.selfOwned && arraysEqual(a.owners, b.owners)
}

function arraysEqual(a: any, b: any) {
    if (a === b) return true;
    if (a == null || b == null) return false;
    if (a.length !== b.length) return false;

    // If you don't care about the order of the elements inside
    // the array, you should sort both arrays here.
    // Please note that calling sort on an array will modify that array.
    // you might want to clone your array first.

    for (var i = 0; i < a.length; ++i) {
        if (a[i] !== b[i]) return false;
    }
    return true;
}
