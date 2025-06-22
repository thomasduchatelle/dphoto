import {Media, MediaWithinADay} from "../language";

export function groupByDay(medias: Media[]): MediaWithinADay[] {
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