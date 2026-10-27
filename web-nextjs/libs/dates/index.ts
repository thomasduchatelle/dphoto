export function toLocaleDateWithDay(date: Date) {
    return new Intl.DateTimeFormat(['en-UK', 'fr-FR'], {
        weekday: 'long',
        day: 'numeric',
        month: 'short',
        year: 'numeric',
    }).format(date);
}

export function toLocaleTimeOfDay(date: Date) {
    return new Intl.DateTimeFormat(['en-UK', 'fr-FR'], {
        hour: 'numeric',
        minute: '2-digit',
    }).format(date);
}
