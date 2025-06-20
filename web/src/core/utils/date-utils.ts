export function toLocaleDateWithDay(date: Date): string {
    return `${date.toLocaleDateString('en-UK', {weekday: 'long', year: 'numeric', month: 'short', day: 'numeric'})}`
}

export function dateTimeToString(date: Date) {
    return `${date.toLocaleDateString()} ${date.toLocaleTimeString()}`
}
