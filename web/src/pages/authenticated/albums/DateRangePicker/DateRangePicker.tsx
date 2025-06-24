import React, {useCallback} from "react";
import {Checkbox, FormControlLabel, TextField} from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
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
    const utcValue = newValue ? newValue.toDate() : null;
    utcValue?.setTime(utcValue.getTime() - utcValue.getTimezoneOffset() * 60 * 1000)
    return utcValue;
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
        <>
            <Grid xs={6}>
                {startAtDayStart ? (
                    <DatePicker
                        label="First day"
                        disabled={disabled}
                        value={startDate ? dayjs(startDate) : null}
                        onChange={onStartChange}
                        renderInput={(params: any) => (
                            <TextField {...params} sx={{width: "100%"}} {...commonDateInputProps} helperText=""/>
                        )}
                    />
                ) : (
                    <DateTimePicker
                        label="First day"
                        disabled={disabled}
                        value={startDate ? dayjs(startDate) : null}
                        onChange={onStartChange}
                        renderInput={(params: any) => (
                            <TextField {...params} sx={{width: "100%"}} {...commonDateInputProps} helperText=""/>
                        )}
                    />
                )}
            </Grid>
            <Grid xs={6}>
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
            </Grid>
            <Grid xs={6}>
                {endAtDayEnd ? (
                    <DatePicker
                        label="Last day"
                        disabled={disabled}
                        value={endDate ? dayjs(endDate) : null}
                        onChange={onEndChange}
                        renderInput={(params: any) => (
                            <TextField {...params} sx={{width: "100%"}} {...commonDateInputProps} />
                        )}
                    />
                ) : (
                    <DateTimePicker
                        label="Last day"
                        disabled={disabled}
                        value={endDate ? dayjs(endDate) : null}
                        onChange={onEndChange}
                        renderInput={(params: any) => (
                            <TextField {...params} sx={{width: "100%"}} {...commonDateInputProps} />
                        )}
                    />
                )}
            </Grid>
            <Grid xs={6}>
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
            </Grid>
        </>
    );
};
