import React from "react";
import {Checkbox, FormControlLabel, TextField} from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import {DatePicker, DateTimePicker} from "@mui/x-date-pickers";

interface DateRangePickerProps {
    startDate: Date;
    endDate: Date;
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

    return (
        <>
            <Grid xs={6}>
                {startAtDayStart ? (
                    <DatePicker
                        label="First day"
                        disabled={disabled}
                        value={startDate}
                        onChange={(newValue: Date | null) => onStartDateChange(newValue)}
                        renderInput={(params: any) => (
                            <TextField {...params} sx={{width: "100%"}} {...commonDateInputProps} helperText=""/>
                        )}
                    />
                ) : (
                    <DateTimePicker
                        label="First day"
                        disabled={disabled}
                        value={startDate}
                        onChange={(newValue: Date | null) => onStartDateChange(newValue)}
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
                        value={endDate}
                        onChange={(newValue: Date | null) => onEndDateChange(newValue)} // Receive Date object
                        renderInput={(params: any) => (
                            <TextField {...params} sx={{width: "100%"}} {...commonDateInputProps} />
                        )}
                    />
                ) : (
                    <DateTimePicker
                        label="Last day"
                        disabled={disabled}
                        value={endDate}
                        onChange={(newValue: Date | null) => onEndDateChange(newValue)} // Receive Date object
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
