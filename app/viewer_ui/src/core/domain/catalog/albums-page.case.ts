import axios, {AxiosError, AxiosInstance} from "axios";
import {AuthenticatedUser} from "../security";
import {Album, AlbumId, CacheAdapter, Media, MediaType, WebAdapter} from "./albums.domain";

interface RestAlbum {
  owner: string
  folderName: string
  name: string
  start: Date
  end: Date
  totalCount: number
}

interface RestMedia {
  id: string
  type: string
  time: string
  path: string
  source: string
}

function numberOfDays(start: Date, end: Date) {
  if (!start || !end) {
    return 1
  }

  return Math.ceil(Math.abs(end.getTime() - start.getTime()) / (1000 * 3600 * 24)) ?? 1;
}

export class AlbumsPageCase {
  constructor(private axios: AxiosInstance,
              private accessToken: string,
              private cacheAdapter: CacheAdapter,
              private webAdapter: WebAdapter) {
  }

  public redirectToDefaultAlbum(loggedUser: AuthenticatedUser): Promise<void> {
    return this.fetchAlbumsWithCache(loggedUser.email).then(albums => {
      if (albums.length === 0) {
        return this.webAdapter.renderNoAlbums()
      } else {
        return this.webAdapter.redirectToAlbum(albums[0].albumId);
      }
    })
  }

  public refreshPage(loggedUser: AuthenticatedUser, albumId: AlbumId): Promise<void> {
    return this.fetchAlbumsWithCache(loggedUser.email)
      .then(albums => {
        const selectedAlbum = albums.find(a => a.albumId.owner === albumId.owner && a.albumId.folderName === albumId.folderName)
        if (!selectedAlbum) {
          return this.webAdapter.renderAlbumNotPresent(albums)
        }

        return this.fetchMedias(albumId)
          .then(medias => this.webAdapter.renderAlbumsWithMedia(albums, selectedAlbum, medias))
          .then()
      })
  }

  private fetchAlbums(owner: string): Promise<Album[]> {
    return this.axios.get<RestAlbum[]>(`/api/v1/owners/${owner}/albums`)
      .then(resp => {
        const maxTemperature = resp.data.map(a => a.totalCount / numberOfDays(new Date(a.start), new Date(a.end))).reduce(function (p, v) {
          return (p > v ? p : v);
        })

        return resp.data.map(album => {
          const temperature = album.totalCount / numberOfDays(new Date(album.start), new Date(album.end));
          return {
            albumId: {owner: album.owner, folderName: album.folderName.replace(/^\//, "")},
            name: album.name,
            start: new Date(album.start),
            end: new Date(album.end),
            totalCount: album.totalCount,
            temperature: temperature,
            relativeTemperature: temperature / maxTemperature,
          }
        }).sort((a, b) => b.start.getTime() - a.start.getTime());
      });
  }

  public selectAlbum = (albumId: AlbumId): void => {
    this.webAdapter.redirectToAlbum(albumId)
  }

  private fetchAlbumsWithCache(loggedEmail: string): Promise<Album[]> {
    const [owner, cachedAlbums] = this.cacheAdapter.getCachedAlbums()

    return owner && owner === loggedEmail ?
      Promise.resolve(cachedAlbums) :
      this.fetchAlbums(loggedEmail)
        .then(albums => {
          this.cacheAdapter.cacheAlbums(loggedEmail, albums)
          return albums
        });
  }

  private fetchMedias(albumId: AlbumId): Promise<Media[]> {
    return this.axios.get<RestMedia[]>(`/api/v1/owners/${albumId.owner}/albums/${albumId.folderName}/medias`)
      .then(resp => resp.data as RestMedia[])
      .catch((err: AxiosError | Error) => {
        if (axios.isAxiosError(err) && err.response?.status === 404) {
          return []
        }
        return Promise.reject<RestMedia[]>(err)
      })
      .then(data => {
        return data.map(media => ({
          ...media,
          type: convertToType(media.type),
          time: new Date(media.time),
          path: `${media.path}?access_token=${this.accessToken}`
        }))
      })
  }
}

function convertToType(type: string): MediaType {
  if (!type) {
    return MediaType.OTHER
  }

  switch (type.toLowerCase()) {
    case 'image':
      return MediaType.IMAGE
    case 'video':
      return MediaType.VIDEO
    default:
      return MediaType.OTHER
  }
}