export function toLocaleDateWithDay(date: Date) {
    return new Intl.DateTimeFormat(['en-GB', 'fr-FR'], {
        weekday: 'long',
        day: 'numeric',
        month: 'short',
        year: 'numeric',
    }).format(date);
}

export function toLocaleTimeOfDay(date: Date) {
    return new Intl.DateTimeFormat(['en-GB', 'fr-FR'], {
        hour: '2-digit',
        minute: '2-digit',
        hourCycle: 'h23',
    }).format(date);
}
