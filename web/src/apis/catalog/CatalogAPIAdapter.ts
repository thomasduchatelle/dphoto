import {Album, AlbumId, CatalogAPI, Media, MediaType, OwnerDetails, UserDetails} from "../../core/catalog";
import axios, {AxiosError, AxiosInstance} from "axios";
import {AccessTokenHolder} from "../../core/application";

interface RestAlbum {
    owner: string
    folderName: string
    name: string
    start: Date
    end: Date
    totalCount: number
    sharedWith?: string[]
    directlyOwned?: boolean
}

interface RestMedia {
    id: string
    type: string
    filename: string
    time: string
    source: string
}

interface RestUserDetails {
    name: string
    email: string
    picture?: string
}

interface RestOwnerDetails {
    id: string
    name?: string
    users: RestUserDetails[]
}

export class CatalogAPIAdapter implements CatalogAPI {
    constructor(
        private readonly authenticatedAxios: AxiosInstance,
        private readonly accessTokenHolder: AccessTokenHolder,
    ) {
    }

    public fetchAlbums(email: string): Promise<Album[]> {
        return this.authenticatedAxios.get<RestAlbum[]>('/api/v1/albums')
            .then(resp => {
                const albums = resp.data;

                return Promise.allSettled([
                    this.findOwnerDetails(new Set<string>(albums.filter(a => !a.directlyOwned).map(a => a.owner))),
                    this.findUserDetails(new Set<string>(albums.flatMap(a => a.sharedWith ?? []))),
                ]).then(([ownersResp, usersResp]) => {
                    const owners = ownersResp.status === "fulfilled" ? ownersResp.value.reduce(
                        (map, owner) => {
                            map.set(owner.id, owner)
                            return map
                        },
                        new Map<string, OwnerDetails>()
                    ) : new Map<string, OwnerDetails>()

                    const users = usersResp.status === "fulfilled" ? usersResp.value.reduce(
                        (map, user) => {
                            map.set(user.email, user)
                            return map
                        },
                        new Map<string, UserDetails>()
                    ) : new Map<string, UserDetails>()

                    return {albums, owners, users}
                })
            })
            .then(({albums, owners, users}) => {

                const maxTemperature = albums.map(a => a.totalCount / numberOfDays(new Date(a.start), new Date(a.end))).reduce(function (p, v) {
                    return (p > v ? p : v);
                })

                return albums.map(album => {
                    const temperature = album.totalCount / numberOfDays(new Date(album.start), new Date(album.end));
                    return {
                        albumId: {owner: album.owner, folderName: album.folderName.replace(/^\//, "")},
                        name: album.name,
                        start: new Date(album.start),
                        end: new Date(album.end),
                        totalCount: album.totalCount,
                        temperature: temperature,
                        relativeTemperature: temperature / maxTemperature,
                        sharedWith: album.sharedWith ? album.sharedWith.map(email => users.get(email) ?? {
                            name: email,
                            email: email,
                        }) : [],
                        ownedBy: album.directlyOwned ? undefined : owners.get(album.owner) ?? {
                            name: album.owner,
                            users: [],
                        }
                    }
                }).sort((a, b) => b.start.getTime() - a.start.getTime());
            });
    }

    private findOwnerDetails(owners: Set<string>): Promise<RestOwnerDetails[]> {
        if (!owners) {
            return Promise.resolve([])
        }

        return this.authenticatedAxios
            .get<RestOwnerDetails[]>(`/api/v1/owners?ids=${[...owners.values()].join(',')}`)
            .then(resp => resp.data);
    }

    private findUserDetails(emails: Set<string>): Promise<RestUserDetails[]> {
        if (!emails) {
            return Promise.resolve([])
        }

        return this.authenticatedAxios
            .get<RestUserDetails[]>(`/api/v1/users?emails=${[...emails.values()].join(',')}`)
            .then(resp => resp.data);
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