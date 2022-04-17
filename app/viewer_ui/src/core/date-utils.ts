const days = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];

export function toLocaleDateWithDay(date: Date) {
  return `${days[date.getDay()]} ${date.toLocaleDateString()}`
}

export function dateTimeToString(date: Date) {
  return `${date.toLocaleDateString()} ${date.toLocaleTimeString()}`
}
