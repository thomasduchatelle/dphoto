import {DateRangeState} from "../language";

export function isRoundTime(date: Date): boolean {
    return date.getUTCMinutes() % 30 === 0 && date.getUTCSeconds() === 0 && date.getUTCMilliseconds() === 0;
}

export interface DateRangeValidation {
    areDatesValid: boolean;
    isDateRangeValid: boolean;
    dateRangeError?: string;
}

export const datesMustBeSetError: DateRangeValidation = {
    areDatesValid: false,
    isDateRangeValid: false,
    dateRangeError: "Both start and end dates must be set",
};

export const startDateMustBeBeforeEndDateError: DateRangeValidation = {
    areDatesValid: true,
    isDateRangeValid: false,
    dateRangeError: "The end date cannot be before the start date",
};

export const dateRangeIsValid: DateRangeValidation = {
    areDatesValid: true,
    isDateRangeValid: true,
};

export function validateDateRange(dateRange: DateRangeState): DateRangeValidation {
    if (!dateRange.startDate || !dateRange.endDate) {
        return datesMustBeSetError
    }

    if (dateRange.startDate > dateRange.endDate) {
        return startDateMustBeBeforeEndDateError;
    }

    return dateRangeIsValid
}

export function convertToModelStartDate(original: Date, atDayStart: boolean): Date {
    const date = new Date(original);
    if (atDayStart) {
        return new Date(Date.UTC(date.getUTCFullYear(), date.getUTCMonth(), date.getUTCDate()));
    }

    return new Date(Date.UTC(date.getUTCFullYear(), date.getUTCMonth(), date.getUTCDate(), date.getHours(), date.getMinutes()));
}

export function convertToModelEndDate(original: Date, atDayEnd: boolean): Date {
    const date = new Date(original);
    if (atDayEnd) {
        return new Date(Date.UTC(date.getUTCFullYear(), date.getUTCMonth(), date.getUTCDate() + 1));
    }

    if (isRoundTime(date)) {
        return new Date(Date.UTC(date.getUTCFullYear(), date.getUTCMonth(), date.getUTCDate(), date.getHours(), date.getMinutes()));
    }

    return new Date(Date.UTC(date.getUTCFullYear(), date.getUTCMonth(), date.getUTCDate(), date.getHours(), date.getMinutes() + 1));
}

export function convertFromModelToDisplayDate(apiDate: Date, isEndDate: boolean, endAtDayEnd: boolean): Date {
    const date = new Date(apiDate);

    if (isEndDate) {
        if (endAtDayEnd) {
            // API end date is exclusive (next day at 00:00), so subtract 1 day for display
            date.setDate(date.getDate() - 1);
        } else if (!isRoundTime(date)) {
            // API end date has +1 minute for exclusivity, so subtract 1 minute for display
            date.setUTCMinutes(date.getUTCMinutes() - 1);
        }
    }

    return date;
}

export function isDateAtDayStart(date: Date): boolean {
    return date.getUTCHours() === 0 && date.getUTCMinutes() === 0 &&
        date.getUTCSeconds() === 0 && date.getUTCMilliseconds() === 0;
}

export function isDateAtDayEnd(date: Date): boolean {
    return date.getUTCHours() === 0 && date.getUTCMinutes() === 0 &&
        date.getUTCSeconds() === 0 && date.getUTCMilliseconds() === 0;
}

export function setDateToStartOfDay(date: Date): Date {
    const newDate = new Date(date);
    newDate.setUTCHours(0, 0, 0, 0);
    return newDate;
}

export function setDateToEndOfDay(date: Date): Date {
    const newDate = new Date(date);
    newDate.setUTCHours(23, 59, 0, 0);
    return newDate;
}
