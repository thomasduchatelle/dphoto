export function toLocaleDate(date: Date): string {
    return date.toLocaleDateString('en-CA', {year: 'numeric', month: '2-digit', day: '2-digit'});
}

export function toLocaleDateWithDay(date: Date): string {
    return `${date.toLocaleDateString('en-CA', {weekday: 'short', year: 'numeric', month: '2-digit', day: '2-digit'})}`
}

export function dateTimeToString(date: Date) {
    return `${date.toLocaleDateString()} ${date.toLocaleTimeString()}`
}
