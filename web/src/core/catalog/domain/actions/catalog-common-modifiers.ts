import {
    Album,
    AlbumFilterCriterion,
    AlbumFilterEntry,
    albumIsOwnedByCurrentUser,
    albumMatchCriterion,
    CurrentUserInsight,
    OwnerDetails
} from "../catalog-state";

export const ALL_ALBUMS_FILTER_CRITERION: AlbumFilterCriterion = {owners: []}
export const SELF_OWNED_ALBUM_FILTER_CRITERION: AlbumFilterCriterion = {selfOwned: true, owners: []}
export const DEFAULT_ALBUM_FILTER_ENTRY: AlbumFilterEntry = {
    criterion: ALL_ALBUMS_FILTER_CRITERION,
    avatars: [],
    name: "All albums",
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


export function albumFilterAreCriterionEqual(a: AlbumFilterCriterion, b: AlbumFilterCriterion) {
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

export function refreshFilters(currentUser: CurrentUserInsight, currentAlbumFilterEntry: AlbumFilterEntry, allAlbums: Album[]) {
    const albumFilterOptions = generateAlbumFilterOptions(currentUser, allAlbums);
    const albumFilter = albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, currentAlbumFilterEntry.criterion)) ?? DEFAULT_ALBUM_FILTER_ENTRY
    const albums = allAlbums.filter(albumMatchCriterion(albumFilter.criterion))

    return {albumFilterOptions, albumFilter, albums};
}

export function filteredListOfAlbums({allAlbums, albumFilter: {criterion}}: {
    allAlbums: Album[]
    albumFilter: AlbumFilterEntry,
}) {
    return allAlbums.filter(albumMatchCriterion(criterion));
}