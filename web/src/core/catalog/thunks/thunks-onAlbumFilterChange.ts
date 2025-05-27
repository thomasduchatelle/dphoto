import {Album, AlbumFilterCriterion, albumIdEquals, albumMatchCriterion, catalogActions, CatalogViewerAction, CatalogViewerState} from "../domain";
import {ThunkDeclaration} from "../../thunk-engine";

import {CatalogFactoryArgs} from "./catalog-factory-args";

export interface AlbumFilterHandlerState {
    selectedAlbum?: Album
    allAlbums: Album[]
}

export function onAlbumFilterFunction(dispatch: (action: CatalogViewerAction) => void, partialState: AlbumFilterHandlerState, criterion: AlbumFilterCriterion) {
    const match = albumMatchCriterion(criterion);
    if (partialState.selectedAlbum && match(partialState.selectedAlbum)) {
        dispatch(catalogActions.albumsFilteredAction({criterion: criterion}));
        return;
    }
    const nextSelectedAlbumId = partialState.allAlbums.find(album => match(album))?.albumId;
    dispatch(catalogActions.albumsFilteredAction({criterion: criterion, redirectTo: nextSelectedAlbumId}));
}

export const onAlbumFilterChangeDeclaration: ThunkDeclaration<
    CatalogViewerState,
    AlbumFilterHandlerState,
    (criterion: AlbumFilterCriterion) => void,
    CatalogFactoryArgs
> = {
    factory: ({dispatch, partialState}) => {
        return onAlbumFilterFunction.bind(null, dispatch, partialState);
    },
    selector: ({mediasLoadedFromAlbumId, loadingMediasFor, allAlbums}: CatalogViewerState): AlbumFilterHandlerState => {
        const albumId = loadingMediasFor || mediasLoadedFromAlbumId;
        const selectedAlbum = allAlbums.find(album => albumId && albumIdEquals(albumId, album.albumId))

        return {
            selectedAlbum,
            allAlbums,
        }
    },
};
