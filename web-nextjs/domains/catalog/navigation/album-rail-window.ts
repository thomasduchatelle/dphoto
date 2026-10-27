import {Album} from "../language";

const RAIL_WINDOW_SIZE = 7;
const RAIL_NEIGHBOURS_EACH_SIDE = 3;

export function albumRailWindow(albums: Album[], displayedIndex: number): Album[] {
    if (albums.length <= RAIL_WINDOW_SIZE || displayedIndex < 0) {
        return albums;
    }

    let start = displayedIndex - RAIL_NEIGHBOURS_EACH_SIDE;
    if (start < 0) {
        start = 0;
    }
    if (start + RAIL_WINDOW_SIZE > albums.length) {
        start = albums.length - RAIL_WINDOW_SIZE;
    }

    return albums.slice(start, start + RAIL_WINDOW_SIZE);
}
