import {AlbumId, albumIdEquals, CatalogViewerState, ShareDialog, ShareError, Sharing, UserDetails} from "../language";
import {filteredListOfAlbums} from "../navigation";

export function sortSharings(sharings: Sharing[]): Sharing[] {
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

export function getSharingSuggestions(
    allAlbums: readonly { sharedWith: Sharing[] }[],
    currentSharedWith: readonly Sharing[],
    extraDetails?: UserDetails[]
): UserDetails[] {
    const knownUsersMap = new Map<string, { user: UserDetails; count: number }>();
    if (extraDetails) {
        for (const user of extraDetails) {
            knownUsersMap.set(user.email, {user, count: 0});
        }
    }

    for (const album of allAlbums) {
        for (const sharing of album.sharedWith) {
            const email = sharing.user.email;
            if (!knownUsersMap.has(email)) {
                knownUsersMap.set(email, {user: sharing.user, count: 1});
            } else {
                knownUsersMap.get(email)!.count += 1;
            }
        }
    }
    const currentEmails = new Set(currentSharedWith.map(s => s.user.email));

    return Array.from(knownUsersMap.values())
        .filter(entry => !currentEmails.has(entry.user.email))
        .sort((a, b) => {
            if (b.count !== a.count) return b.count - a.count;
            const nameA = a.user.name?.trim() || "";
            const nameB = b.user.name?.trim() || "";
            return nameA.localeCompare(nameB);
        })
        .map(entry => entry.user);
}

export function withOpenShareDialog(state: CatalogViewerState, albumId: AlbumId, extra: Partial<ShareDialog> = {}): CatalogViewerState {
    const album = state.allAlbums.find(a => albumIdEquals(a.albumId, albumId));
    return {
        ...state,
        dialog: album ? {
            type: "ShareDialog",
            sharedAlbumId: albumId,
            sharedWith: sortSharings([...album.sharedWith]),
            ...extra,
            suggestions: getSharingSuggestions(state.allAlbums, album.sharedWith, extra.suggestions ?? []),
        } : undefined,
    }
}

export function moveSharedWithToSuggestion(
    {allAlbums, albums: _, ...current}: CatalogViewerState,
    {sharedWith: previousSharedWith, suggestions: previousSuggestions, error: _2, ...shareDialog}: ShareDialog,
    email: string,
    error ?: ShareError
): CatalogViewerState {
    const restoredAllAlbums = allAlbums.map(album => {
        if (albumIdEquals(album.albumId, shareDialog.sharedAlbumId)) {
            return {
                ...album,
                sharedWith: album.sharedWith.filter(s => s.user.email !== email),
            };
        }

        return album;
    });

    return withOpenShareDialog(
        {
            ...current,
            allAlbums: restoredAllAlbums,
            albums: filteredListOfAlbums({...current, allAlbums: restoredAllAlbums}),
        },
        shareDialog.sharedAlbumId,
        {
            suggestions: previousSuggestions.concat(previousSharedWith.filter(({user}) => user.email === email).map(({user}) => user)),
            error,
        }
    );
}

export function moveSuggestionToSharedWith(
    {allAlbums, albums: _, ...current}: CatalogViewerState,
    {sharedWith: previousSharedWith, suggestions: previousSuggestions, error: _2, ...shareDialog}: ShareDialog,
    user: UserDetails,
    error?: ShareError,
): CatalogViewerState {
    const restoredAllAlbums = allAlbums.map(album => {
        if (albumIdEquals(album.albumId, shareDialog.sharedAlbumId) && !album.sharedWith.some(s => s.user.email === user.email)) {
            const sharedWith = sortSharings([...album.sharedWith, {user}]);
            return {
                ...album,
                sharedWith,
            };
        }
        return album;
    });

    return withOpenShareDialog(
        {
            ...current,
            allAlbums: restoredAllAlbums,
            albums: filteredListOfAlbums({...current, allAlbums: restoredAllAlbums}),
        },
        shareDialog.sharedAlbumId,
        {
            suggestions: previousSuggestions.filter(s => s.email !== user.email),
            error,
        }
    );
}
