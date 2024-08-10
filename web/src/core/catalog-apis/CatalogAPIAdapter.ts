import {Album, AlbumId, CatalogAPI, Media, MediaType, OwnerDetails, SharingType, UserDetails} from "../catalog-react";
import axios, {AxiosError, AxiosInstance} from "axios";
import {AccessTokenHolder} from "../application";
import {SharingAPI} from "../../pages/authenticated/albums/share-controller";

interface RestAlbum {
    owner: string
    folderName: string
    name: string
    start: Date
    end: Date
    totalCount: number
    sharedWith?: Record<string, string>
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

export class CatalogAPIAdapter implements CatalogAPI, SharingAPI {
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
                    this.findUserDetails(new Set<string>(albums.flatMap(a => Object.entries(a.sharedWith ?? {}).map(([email, role]) => email)))),
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
                        sharedWith: Object.entries(album.sharedWith ?? {}).map(([email, role]) => ({
                            user: users.get(email) ?? {
                                name: email,
                                email: email,
                            },
                            role: castAsSharingType(role, SharingType.visitor),
                        })),
                        ownedBy: album.directlyOwned ? undefined : owners.get(album.owner) ?? {
                            name: album.owner,
                            users: [],
                        }
                    }
                }).sort((a, b) => b.start.getTime() - a.start.getTime());
            });
    }

    private findOwnerDetails(owners: Set<string>): Promise<RestOwnerDetails[]> {
        if (owners.size === 0) {
            return Promise.resolve([])
        }

        return this.authenticatedAxios
            .get<RestOwnerDetails[]>(`/api/v1/owners`, {
                params: {
                    ids: [...owners.values()].join(','),
                }
            })
            .then(resp => resp.data);
    }

    private findUserDetails(emails: Set<string>): Promise<RestUserDetails[]> {
        if (emails.size === 0) {
            return Promise.resolve([])
        }

        return this.authenticatedAxios
            .get<RestUserDetails[]>(`/api/v1/users`, {
                params: {
                    emails: [...emails.values()].join(','),
                }
            })
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

    public grantAccessToAlbum(albumId: AlbumId, email: string, role: SharingType): Promise<void> {
        return this.authenticatedAxios
            .put(`/api/v1/owners/${albumId.owner}/albums/${albumId.folderName}/shares/${email}`, {
                role,
            })
    }

    public loadUserDetails(email: string): Promise<UserDetails> {
        return this.findUserDetails(new Set<string>([email]))
            .then(details => {
                if (details && details.length > 0) {
                    return Promise.resolve({...details[0]})
                }

                return Promise.reject("user details not found.")
            })
    }

    public revokeSharingAlbum(albumId: AlbumId, email: string): Promise<void> {
        return this.authenticatedAxios
            .delete(`/api/v1/owners/${albumId.owner}/albums/${albumId.folderName}/shares/${email}`)
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

function castAsSharingType(value: string, defaultValue: SharingType): SharingType {
    return (SharingType as any)[value] ?? defaultValue
}
