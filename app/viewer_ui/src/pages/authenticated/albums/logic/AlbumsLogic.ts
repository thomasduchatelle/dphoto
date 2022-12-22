import axios, {AxiosError} from "axios";
import {MustBeAuthenticated} from "../../../../core/application";
import {Cache} from "../../../../core/common";
import {Album, AlbumId, AlbumsLogicCache, Media, MediaType, MediaWithinADay, WebAdapter} from "./albums.domain";

export const mobileBreakpoint = 1200;

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
  filename: string
  time: string
  source: string
}

function numberOfDays(start: Date, end: Date) {
  if (!start || !end) {
    return 1
  }

  return Math.ceil(Math.abs(end.getTime() - start.getTime()) / (1000 * 3600 * 24)) ?? 1;
}

export class AlbumsLogic {
  constructor(
    private mustBeAuthenticated: MustBeAuthenticated,
    private cache: Cache<AlbumsLogicCache>,
    private width: number,
    private webAdapter: WebAdapter,
    private mobileMode: boolean = false
  ) {
    this.mobileMode = width < mobileBreakpoint
  }

  public loadAlbumsPage = (): Promise<void> => {
    return this.fetchAlbumsWithCache(this.mustBeAuthenticated.loggedUser.email).then(albums => {
      if (albums.length === 0) {
        return this.webAdapter.renderNoAlbums()

      } else if (this.mobileMode) {
        return this.webAdapter.renderAlbumsList(albums)

      } else {
        return this.webAdapter.redirectToAlbum(albums[0].albumId);
      }
    })
  }

  public loadMediasPage = (albumId: AlbumId): Promise<void> => {
    return this.fetchAlbumsWithCache(this.mustBeAuthenticated.loggedUser.email)
      .then(albums => {
        const selectedAlbum = albums.find(a => a.albumId.owner === albumId.owner && a.albumId.folderName === albumId.folderName)
        if (!selectedAlbum) {
          return this.webAdapter.renderAlbumNotPresent(albums, albumId)
        }

        return this.fetchMedias(albumId)
          .then(medias => this.groupByDay(medias))
          .then(medias => this.webAdapter.renderAlbumsWithMedia(albums, selectedAlbum, medias))
          .then()
      })
  }

  private groupByDay(medias: Media[]): MediaWithinADay[] {
    let result: MediaWithinADay[] = []

    medias.forEach(m => {
      const beginning = new Date(m.time)
      beginning.setHours(0, 0, 0, 0)

      if (result.length > 0 && result[0].day.getTime() === beginning.getTime()) {
        result[0].medias.push(m)
      } else {
        result = [{
          day: beginning,
          medias: [m],
        }, ...result]
      }
    })

    result.reverse()
    return result
  }

  private fetchAlbums(owner: string): Promise<Album[]> {
    return this.mustBeAuthenticated.authenticatedAxios.get<RestAlbum[]>(`/api/v1/albums`)
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
    const {owner, albums} = this.cache.current

    if (owner && owner === loggedEmail) {
      return Promise.resolve(albums)
    }

    return this.fetchAlbums(loggedEmail).then(albums => {
      this.cache.current = {owner: loggedEmail, albums}
      return albums
    });
  }

  private fetchMedias(albumId: AlbumId): Promise<Media[]> {
    return this.mustBeAuthenticated.authenticatedAxios
      .get<RestMedia[]>(`/api/v1/owners/${albumId.owner}/albums/${albumId.folderName}/medias`)
      .then(resp => resp.data as RestMedia[])
      .catch((err: AxiosError | Error) => {
        if (axios.isAxiosError(err) && err.response?.status === 404) {
          return []
        }
        return Promise.reject<RestMedia[]>(err)
      })
      .then(data => {
        return data.map((media): Media => ({
          ...media,
          type: convertToType(media.type),
          time: new Date(media.time),
          uiRelativePath: `${media.id}/${media.filename}`,
          contentPath: `/api/v1/owners/${albumId.owner}/medias/${media.id}/${media.filename}?access_token=${this.mustBeAuthenticated.accessToken}`,
      })).sort((a, b) => b.time.getTime() - a.time.getTime())
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