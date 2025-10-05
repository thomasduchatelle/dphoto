'use client';

import React, {useCallback} from "react";
import {Checkbox, FormControlLabel, TextField, Box} from "@mui/material";
import {DatePicker, DateTimePicker} from "@mui/x-date-pickers";
import dayjs, {Dayjs} from "dayjs";
// import utc from 'dayjs/plugin/utc';

// dayjs.extend(utc)

interface DateRangePickerProps {
    startDate: Date | null;
    endDate: Date | null;
    startAtDayStart: boolean;
    endAtDayEnd: boolean;
    onStartDateChange: (date: Date | null) => void;
    onEndDateChange: (date: Date | null) => void;
    onStartsAtStartOfTheDayChange: (checked: boolean) => void;
    onEndsAtEndOfTheDayChange: (checked: boolean) => void;
    disabled?: boolean;
    dateError?: boolean;
    dateHelperText?: string;
}


function toUTCDate(newValue: Dayjs | null) {
    if (!newValue) {
        return null;
    }

    const utcValue = newValue.toDate();
    utcValue.setTime(utcValue.getTime() - utcValue.getTimezoneOffset() * 60 * 1000)
    return utcValue;
}

function fromUTCDate(utcDate: Date | null): Dayjs | null {
    if (!utcDate) {
        return null;
    }

    const hacked = new Date(utcDate)
    hacked.setTime(hacked.getTime() + hacked.getTimezoneOffset() * 60 * 1000)

    return dayjs(hacked);
}

export const DateRangePicker: React.FC<DateRangePickerProps> = ({
                                                                    startDate,
                                                                    endDate,
                                                                    startAtDayStart,
                                                                    endAtDayEnd,
                                                                    onStartDateChange,
                                                                    onEndDateChange,
                                                                    onStartsAtStartOfTheDayChange,
                                                                    onEndsAtEndOfTheDayChange,
                                                                    disabled = false,
                                                                    dateError = false,
                                                                    dateHelperText = "",
                                                                }) => {
    const commonDateInputProps = {
        error: dateError,
        helperText: dateHelperText,
    };

    const onStartChange = useCallback((newValue: Dayjs | null) => onStartDateChange(toUTCDate(newValue)), [onStartDateChange]);
    const onEndChange = useCallback((newValue: Dayjs | null) => onEndDateChange(toUTCDate(newValue)), [onEndDateChange]);

    return (
        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 2 }}>
            <Box sx={{ flexBasis: { xs: '100%', sm: '48%' } }}>
                {startAtDayStart ? (
                    <DatePicker
                        label="First day"
                        disabled={disabled}
                        value={fromUTCDate(startDate)}
                        onChange={onStartChange}
                        slots={{
                            textField: TextField
                        }}
                        slotProps={{
                            textField: {
                                sx: {width: "100%"},
                                ...commonDateInputProps,
                                helperText: ""
                            }
                        }}
                    />
                ) : (
                    <DateTimePicker
                        label="First day"
                        disabled={disabled}
                        value={fromUTCDate(startDate)}
                        onChange={onStartChange}
                        slots={{
                            textField: TextField
                        }}
                        slotProps={{
                            textField: {
                                sx: {width: "100%"},
                                ...commonDateInputProps,
                                helperText: ""
                            }
                        }}
                    />
                )}
            </Box>
            <Box sx={{ flexBasis: { xs: '100%', sm: '48%' } }}>
                <FormControlLabel
                    control={
                        <Checkbox
                            checked={startAtDayStart}
                            disabled={disabled}
                            onChange={(event: React.ChangeEvent<HTMLInputElement>) =>
                                onStartsAtStartOfTheDayChange(event.target.checked)
                            }
                        />
                    }
                    label="at the start of the day"
                />
            </Box>
            <Box sx={{ flexBasis: { xs: '100%', sm: '48%' } }}>
                {endAtDayEnd ? (
                    <DatePicker
                        label="Last day"
                        disabled={disabled}
                        value={fromUTCDate(endDate)}
                        onChange={onEndChange}
                        slots={{
                            textField: TextField
                        }}
                        slotProps={{
                            textField: {
                                sx: {width: "100%"},
                                ...commonDateInputProps
                            }
                        }}
                    />
                ) : (
                    <DateTimePicker
                        label="Last day"
                        disabled={disabled}
                        value={fromUTCDate(endDate)}
                        onChange={onEndChange}
                        slots={{
                            textField: TextField
                        }}
                        slotProps={{
                            textField: {
                                sx: {width: "100%"},
                                ...commonDateInputProps
                            }
                        }}
                    />
                )}
            </Box>
            <Box sx={{ flexBasis: { xs: '100%', sm: '48%' } }}>
                <FormControlLabel
                    control={
                        <Checkbox
                            checked={endAtDayEnd}
                            disabled={disabled}
                            onChange={(event: React.ChangeEvent<HTMLInputElement>) =>
                                onEndsAtEndOfTheDayChange(event.target.checked)
                            }
                        />
                    }
                    label="at the end of the day"
                />
            </Box>
        </Box>
    );
};
