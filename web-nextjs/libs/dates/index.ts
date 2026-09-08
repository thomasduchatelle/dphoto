export function toLocaleDateWithDay(date: Date) {
    return new Intl.DateTimeFormat(['en-UK', 'fr-FR'], {
        weekday: 'long',
        day: 'numeric',
        month: 'short',
        year: 'numeric',
    }).format(date);
}
