export function isRoundTime(date: Date): boolean {
    const minutes = date.getMinutes();
    return minutes === 0 || minutes === 30;
}
