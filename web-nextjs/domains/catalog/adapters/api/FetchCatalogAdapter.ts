import {Album, AlbumId, CatalogError, Media, MediaType, OwnerDetails, UserDetails} from "../../language";
import {GrantAlbumAccessAPI, RevokeAlbumAccessAPI} from "../../sharing";
import {DeleteAlbumPort, SaveAlbumNamePort, UpdateAlbumDatesPort} from "@/domains/catalog";
import {CreateAlbumPort, CreateAlbumRequest} from "../../album-create/thunk-submitCreateAlbum";

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

export type MasterCatalogAdapter = CreateAlbumPort & GrantAlbumAccessAPI & RevokeAlbumAccessAPI & DeleteAlbumPort & UpdateAlbumDatesPort & SaveAlbumNamePort

export class FetchCatalogAdapter implements MasterCatalogAdapter {
    constructor(
        private readonly accessTokenSupplier: () => Promise<string | undefined>,
        private readonly baseUrlSupplier: () => Promise<string> = async () => "/api/v1",
    ) {
    }

    public async deleteAlbum(albumId: AlbumId): Promise<void> {
        await this.fetchRequest(
            `/owners/${albumId.owner}/albums/${albumId.folderName}`,
            {method: 'DELETE'}
        );
    }

    public async createAlbum(request: CreateAlbumRequest): Promise<AlbumId> {
        return this.fetchRequest<AlbumId>(`/albums`, {
            method: 'POST',
            body: JSON.stringify({
                name: request.name,
                folderName: request.forcedFolderName,
                start: request.start.toISOString(),
                end: request.end.toISOString(),
            })
        });
    }

    public fetchAlbums(): Promise<Album[]> {
        return this.fetchRequest<RestAlbum[]>(`/albums`)
            .then(albums => {
                return Promise.allSettled([
                    this.findOwnerDetails(new Set<string>(albums.filter(a => !a.directlyOwned).map(a => a.owner))),
                    this.findUserDetails(new Set<string>(albums.flatMap(a => Object.entries(a.sharedWith ?? {}).map(([email]) => email)))),
                ]).then(([ownersResp, usersResp]) => {
                    const owners = ownersResp.status === "fulfilled" ? ownersResp.value.reduce(
                        (map, owner) => {
                            map.set(owner.id, {name: owner.name, users: owner.users} as OwnerDetails)
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
                        sharedWith: Object.entries(album.sharedWith ?? {}).map(([email]) => ({
                            user: users.get(email) ?? {
                                name: email,
                                email: email,
                            },
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

        const ids = [...owners.values()].join(',');
        return this.fetchRequest<RestOwnerDetails[]>(`/owners?ids=${encodeURIComponent(ids)}`);
    }

    private findUserDetails(emails: Set<string>): Promise<RestUserDetails[]> {
        if (emails.size === 0) {
            return Promise.resolve([])
        }

        const emailsParam = [...emails.values()].join(',');
        return this.fetchRequest<RestUserDetails[]>(`/users?emails=${encodeURIComponent(emailsParam)}`);
    }

    public fetchMedias(albumId: AlbumId): Promise<Media[]> {
        return this.fetchRequest<RestMedia[]>(
            `/owners/${albumId.owner}/albums/${albumId.folderName}/medias`
        )
            .catch((err: Error) => {
                if (err instanceof CatalogError && err.message.includes('404')) {
                    return []
                }
                return Promise.reject<RestMedia[]>(err)
            })
            .then(data => {
                return data.map((media): Media => ({
                    id: media.id,
                    source: media.source,
                    type: convertToType(media.type),
                    time: new Date(media.time),
                    uiRelativePath: `/albums/${albumId.owner}/${albumId.folderName}/${media.id}/${media.filename}`,
                    contentPath: `/api/v1/owners/${albumId.owner}/medias/${media.id}/${media.filename}`,
                })).sort((a, b) => b.time.getTime() - a.time.getTime())
            })
    }

    public grantAccessToAlbum(albumId: AlbumId, email: string): Promise<void> {
        return this.fetchRequest(
            `/owners/${albumId.owner}/albums/${albumId.folderName}/shares/${email}`,
            {method: 'PUT'}
        );
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
        return this.fetchRequest(
            `/owners/${albumId.owner}/albums/${albumId.folderName}/shares/${email}`,
            {method: 'DELETE'}
        );
    }

    public async updateAlbumDates(albumId: AlbumId, startDate: Date, endDate: Date): Promise<void> {
        await this.fetchRequest(
            `/owners/${albumId.owner}/albums/${albumId.folderName}/dates`,
            {
                method: 'PUT',
                body: JSON.stringify({
                    start: startDate.toISOString(),
                    end: endDate.toISOString(),
                })
            }
        );
    }

    public async renameAlbum(albumId: AlbumId, newName: string, newFolderName?: string): Promise<AlbumId> {
        return this.fetchRequest<AlbumId>(
            `/owners/${albumId.owner}/albums/${albumId.folderName}/name`,
            {
                method: 'PUT',
                body: JSON.stringify({
                    name: newName,
                    folderName: newFolderName,
                })
            }
        );
    }

    private async fetchRequest<T>(url: string, options?: RequestInit): Promise<T> {
        const baseUrl = await this.baseUrlSupplier();
        const accessToken = await this.accessTokenSupplier();

        const defaultOptions: RequestInit = {
            headers: {
                'Content-Type': 'application/json',
                ...(accessToken ? {'Authorization': `Bearer ${accessToken}`} : {}),
            },
            credentials: 'include',
        };

        try {
            let fullUrl = `${baseUrl}${url}`;
            console.log("Requesting:", fullUrl, options);
            const response = await fetch(fullUrl, {...defaultOptions, ...options});

            if (!response.ok) {
                const contentType = response.headers.get('content-type');
                if (contentType && contentType.includes('application/json')) {
                    const errorData = await response.json();
                    if (errorData.code && errorData.message) {
                        throw new CatalogError(errorData.code, errorData.message);
                    }
                }

                const defaultMessage = `'${options?.method ?? 'GET'} ${fullUrl}' failed with status ${response.status} ${response.statusText}`;
                throw new CatalogError('', defaultMessage);
            }

            if (response.status === 204 || response.headers.get('content-length') === '0') {
                return undefined as T;
            }

            return response.json();
        } catch (err) {
            if (err instanceof CatalogError) {
                throw err;
            }
            throw new CatalogError('', `Request failed: ${err instanceof Error ? err.message : 'Unknown error'}`);
        }
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
