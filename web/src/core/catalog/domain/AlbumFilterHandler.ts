import {Album, AlbumFilterCriterion, albumMatchCriterion} from "./catalog-state";
import {AlbumsFilteredAction, MediasLoadedAction} from "./catalog-actions";

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
            this.dispatch({
                type: "AlbumsFilteredAction",
                criterion: criterion,
            })
            return
        }

        const nextSelectedAlbumId = this.partialState.allAlbums.find(album => match(album))?.albumId;
        this.dispatch({
            type: "AlbumsFilteredAction",
            criterion: criterion,
            albumId: nextSelectedAlbumId,
        })
    }
}