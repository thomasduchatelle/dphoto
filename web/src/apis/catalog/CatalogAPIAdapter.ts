import {Album, AlbumId, CatalogAPI, Media, MediaType} from "../../core/catalog";
import axios, {AxiosError, AxiosInstance} from "axios";
import {AccessTokenHolder} from "../../core/application";

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

export class CatalogAPIAdapter implements CatalogAPI {
    constructor(
        private readonly authenticatedAxios: AxiosInstance,
        private readonly accessTokenHolder: AccessTokenHolder,
    ) {
    }

    public fetchAlbums(email: string): Promise<Album[]> {
        return this.authenticatedAxios.get<RestAlbum[]>(`/api/v1/albums`)
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

    public fetchMedias(albumId: AlbumId): Promise<Media[]> {
        return this.authenticatedAxios
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
                    contentPath: `/api/v1/owners/${albumId.owner}/medias/${media.id}/${media.filename}?access_token=${this.accessTokenHolder.getAccessToken()}`,
                })).sort((a, b) => b.time.getTime() - a.time.getTime())
            })
    }
}

function numberOfDays(start: Date, end: Date) {
    if (!start || !end) {
        return 1
    }

    return Math.ceil(Math.abs(end.getTime() - start.getTime()) / (1000 * 3600 * 24)) ?? 1;
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