import {AlbumId} from "./catalog-state";

export interface AlbumTemperatureInput {
    albumId: AlbumId
    start: Date
    end: Date
    totalCount: number
}

export interface AlbumTemperature {
    temperature: number
    relativeTemperature: number
}

export function computeAlbumTemperatures(albums: AlbumTemperatureInput[]): Map<string, AlbumTemperature> {
    const n = albums.length
    const result = new Map<string, AlbumTemperature>()
    if (n === 0) {
        return result
    }

    const withTemperature = albums
        .map(album => ({
            key: albumKey(album.albumId),
            temperature: album.totalCount / numberOfDays(album.start, album.end),
        }))
        .sort((a, b) => a.temperature - b.temperature)

    let i = 0
    while (i < n) {
        let j = i
        while (j < n && withTemperature[j].temperature === withTemperature[i].temperature) {
            j++
        }

        const relativeTemperature = (i + j) / (2 * n)
        for (let k = i; k < j; k++) {
            result.set(withTemperature[k].key, {temperature: withTemperature[k].temperature, relativeTemperature})
        }
        i = j
    }

    return result
}

export function albumKey(albumId: AlbumId): string {
    return `${albumId.owner}|${albumId.folderName}`
}

function numberOfDays(start: Date, end: Date): number {
    return Math.max(1, Math.ceil(Math.abs(end.getTime() - start.getTime()) / (1000 * 3600 * 24)))
}
