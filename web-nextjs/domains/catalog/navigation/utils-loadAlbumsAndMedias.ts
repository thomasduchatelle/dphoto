import {Album, AlbumId, albumIdEquals} from "../language";
import {noAlbumAvailable, NoAlbumAvailable} from "./action-noAlbumAvailable";
import {albumsAndMediasLoaded, AlbumsAndMediasLoaded} from "./action-albumsAndMediasLoaded";
import {mediaLoadFailed, MediaLoadFailed} from "./action-mediaLoadFailed";
import {FetchAlbumsAndMediasPort} from "./thunk-onPageRefresh";

export async function loadAlbumsAndMedias(fetchAlbumsAndMediasPort: FetchAlbumsAndMediasPort, albumId ?: AlbumId): Promise<NoAlbumAvailable | AlbumsAndMediasLoaded | MediaLoadFailed> {
    let albums: Album[] = [];

    try {
        albums = await fetchAlbumsAndMediasPort.fetchAlbums()
        if (!albums || albums.length === 0) {
            return noAlbumAvailable(undefined)
        }
    } catch (e) {
        console.log("loadAlbumsAndMedias > loading albums failed", e)
        return noAlbumAvailable(e as Error);
    }

    const albumIdToLoad = albums.find(a => albumIdEquals(a.albumId, albumId))?.albumId ?? albums[0].albumId;
    try {
        const medias = await fetchAlbumsAndMediasPort.fetchMedias(albumIdToLoad);
        return albumsAndMediasLoaded({
            albums: albums,
            medias: medias,
            mediasFromAlbumId: albumIdToLoad,
            redirectTo: albumIdToLoad,
        });

    } catch (e: any) {
        console.log("loadAlbumsAndMedias > loading medias failed", e)
        return mediaLoadFailed({
            albums: albums,
            displayedAlbumId: albumIdToLoad,
            error: new Error(`failed to load medias of ${albumIdToLoad}`, e),
        })
    }
}