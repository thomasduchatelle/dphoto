import {Album, AlbumFilterCriterion, AlbumId, CatalogViewerState, CreateAlbumListener, SharingType} from "../../catalog";

export interface CatalogHandlers extends CreateAlbumListener {
    onAlbumFilterChange: (criterion: AlbumFilterCriterion) => void
}

export interface ShareHandlers {

    onRevoke(email: string): Promise<void>

    onGrant(email: string, role: SharingType): Promise<void>

    openSharingModal(album: Album): void

    onClose(): void
}

export interface CatalogViewerStateWithDispatch {
    state: CatalogViewerState
    selectedAlbumId?: AlbumId // state managed from the URL
    handlers: CatalogHandlers & ShareHandlers
}