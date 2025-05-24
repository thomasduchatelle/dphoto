import {Album, AlbumFilterCriterion, albumMatchCriterion} from "./catalog-state";
import {AlbumsFilteredAction, MediasLoadedAction, catalogActions} from "./catalog-reducer-v2";

export interface AlbumFilterHandlerState {
    selectedAlbum?: Album
    allAlbums: Album[]
}

export type AlbumFilterHandlerDispatch = (action: AlbumsFilteredAction | MediasLoadedAction) => void

export class AlbumFilterHandler {
    constructor(
        private readonly dispatch: AlbumFilterHandlerDispatch,
        private readonly partialState: AlbumFilterHandlerState,
    ) {
    }

    public onAlbumFilter = (criterion: AlbumFilterCriterion,) => {
        const match = albumMatchCriterion(criterion);
        if (this.partialState.selectedAlbum && match(this.partialState.selectedAlbum)) {
            this.dispatch(catalogActions.albumsFilteredAction({criterion: criterion}))
            return
        }

        const nextSelectedAlbumId = this.partialState.allAlbums.find(album => match(album))?.albumId;
        this.dispatch(catalogActions.albumsFilteredAction({criterion: criterion, redirectTo: nextSelectedAlbumId}))
    }
}
