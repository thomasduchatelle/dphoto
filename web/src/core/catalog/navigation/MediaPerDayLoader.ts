import {AlbumId, Media, MediaWithinADay} from "../language";


export interface FetchAlbumMediasPort {
    fetchMedias(albumId: AlbumId): Promise<Media[]>
}

export class MediaPerDayLoader {

    constructor(
        private readonly fetchAlbumMediasPort: FetchAlbumMediasPort,
    ) {
    }

    public async loadMedias(albumId: AlbumId): Promise<MediaWithinADay[]> {
        const medias = await this.fetchAlbumMediasPort.fetchMedias(albumId)
        return this.groupByDay(medias)
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
}
