export interface CatalogViewerState {
    currentUser: CurrentUserInsight // immutable - the whole context is unmounted and remounted when the user changes

    albumNotFound: boolean
    allAlbums: Album[]
    albumFilterOptions: AlbumFilterEntry[],
    albumFilter: AlbumFilterEntry,
    // albums is the list derived from 'albumFilter' and 'allAlbums'
    albums: Album[]
    mediasLoadedFromAlbumId?: AlbumId
    medias: MediaWithinADay[]
    error?: Error
    loadingMediasFor?: AlbumId
    albumsLoaded: boolean
    mediasLoaded: boolean
    dialog?: CatalogDialog
}

export interface CurrentUserInsight {
    picture?: string
}

export interface DateRangeState {
    startDate: Date | null;
    endDate: Date | null;
    startAtDayStart: boolean;
    endAtDayEnd: boolean;
}

export enum MediaType {
    IMAGE,
    VIDEO,
    OTHER
}

export type MediaId = string

export interface Media {
    id: MediaId
    type: MediaType
    time: Date
    uiRelativePath: string // uiRelativePath is the WEB UI internal link (from the album)
    contentPath: string
    source: string
}

export interface MediaWithinADay {
    day: Date
    medias: Media[]
}

export type Owner = string

export interface AlbumId {
    owner: Owner
    folderName: string
}

export interface OwnerDetails {
    name: string
    users: UserDetails[]
}

export interface Sharing {
    user: UserDetails
}

export interface UserDetails {
    name: string
    email: string
    picture?: string
}

export interface Album {
    albumId: AlbumId
    name: string
    start: Date
    end: Date
    totalCount: number
    temperature: number // number of media per day
    relativeTemperature: number
    ownedBy?: OwnerDetails // only present when not owned by current user TODO should be present when owned by user or his picture won't be available.
    sharedWith: Sharing[]
}

export interface AlbumFilterCriterion {
    owners: Owner[] // Empty with selfOwned=false means all albums user has access to
    selfOwned?: boolean // Owned by the current user
}

export interface AlbumFilterEntry {
    criterion: AlbumFilterCriterion
    avatars: string[]
    name: string
}

export interface ShareError {
    type: "grant" | "revoke"
    message: string
    email: string
}

export interface ShareModal {
    sharedAlbumId: AlbumId
    sharedWith: Sharing[]
    suggestions: UserDetails[]
    error?: ShareError
}

export interface ShareDialog {
    type: "ShareDialog"
    sharedAlbumId: AlbumId
    sharedWith: Sharing[]
    suggestions: UserDetails[]
    error?: ShareError
}

export function albumIsOwnedByCurrentUser(album: Album) {
    return album.ownedBy === undefined;
}

export function albumMatchCriterion(criterion: AlbumFilterCriterion): (album: Album) => boolean {
    return album => {
        if (criterion.selfOwned) {
            return albumIsOwnedByCurrentUser(album)
        } else {
            return criterion.owners.length === 0 || criterion.owners.includes(album.albumId.owner)
        }
    };
}

export type RedirectToAlbumIdPayload = {
    redirectTo?: AlbumId
}

export function isRedirectToAlbumIdPayload(arg: any): arg is RedirectToAlbumIdPayload {
    return arg && arg.redirectTo
}

export interface DeleteDialog {
    type: "DeleteDialog";
    deletableAlbums: Album[]
    initialSelectedAlbumId?: AlbumId
    isLoading: boolean
    error?: string
}

export interface CreateDialog extends DateRangeState {
    type: "CreateDialog"
    name: string
    forceFolderName: string
    withCustomFolderName: boolean
    isLoading: boolean
    error?: string
}

export interface EditDatesDialog extends DateRangeState {
    type: "EditDatesDialog"
    albumId: AlbumId
    albumName: string
    isLoading: boolean
    error?: string
}

export type CatalogDialog = CreateDialog | EditDatesDialog | DeleteDialog | ShareDialog;

export function isEditDatesDialog(dialog: CatalogDialog | undefined): dialog is EditDatesDialog {
    return dialog?.type === "EditDatesDialog";
}

export function isDeleteDialog(dialog: CatalogDialog | undefined): dialog is DeleteDialog {
    return dialog?.type === "DeleteDialog";
}

export function isCreateDialog(dialog: CatalogDialog | undefined): dialog is CreateDialog {
    return dialog?.type === "CreateDialog";
}

export function isShareDialog(dialog: CatalogDialog | undefined): dialog is ShareDialog {
    return dialog?.type === "ShareDialog";
}
