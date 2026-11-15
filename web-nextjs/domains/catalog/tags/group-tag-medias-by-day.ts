import {TagMedia, TagMediaWithinADay} from "./language";

export function groupTagMediasByDay(medias: TagMedia[]): TagMediaWithinADay[] {
    const sorted = [...medias].sort((a, b) => b.time.getTime() - a.time.getTime());

    const result: TagMediaWithinADay[] = [];

    sorted.forEach(m => {
        const beginning = new Date(m.time);
        beginning.setHours(0, 0, 0, 0);

        const lastDay = result[result.length - 1];
        if (lastDay && lastDay.day.getTime() === beginning.getTime()) {
            lastDay.medias.push(m);
        } else {
            result.push({day: beginning, medias: [m]});
        }
    });

    return result;
}
