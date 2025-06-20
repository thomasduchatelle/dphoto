import {Album, AlbumFilterCriterion, albumIdEquals, albumMatchCriterion, CatalogViewerState} from "../language";
import {AlbumsFiltered, albumsFiltered} from "./action-albumsFiltered";
import {ThunkDeclaration} from "../../thunk-engine";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";

export interface AlbumFilterHandlerState {
    selectedAlbum?: Album
    allAlbums: Album[]
}

export function onAlbumFilterFunction(dispatch: (action: AlbumsFiltered) => void, partialState: AlbumFilterHandlerState, criterion: AlbumFilterCriterion) {
    const match = albumMatchCriterion(criterion);
    if (partialState.selectedAlbum && match(partialState.selectedAlbum)) {
        dispatch(albumsFiltered({criterion: criterion}));
        return;
    }
    const nextSelectedAlbumId = partialState.allAlbums.find(album => match(album))?.albumId;
    dispatch(albumsFiltered({criterion: criterion, redirectTo: nextSelectedAlbumId}));
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
