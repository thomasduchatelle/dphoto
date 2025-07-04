import React from "react";
import {Alert, Box, Button, Dialog, DialogActions, DialogContent, DialogTitle, IconButton, LinearProgress, useMediaQuery, useTheme} from "@mui/material";
import {Close} from "@mui/icons-material";
import Grid from "@mui/material/Unstable_Grid2";
import {DateRangePicker} from "../DateRangePicker";

interface EditDatesDialogProps {
    isOpen: boolean;
    albumName?: string;
    startDate: Date | null;
    endDate: Date | null;
    startAtDayStart?: boolean;
    endAtDayEnd?: boolean;
    isLoading: boolean;
    errorCode?: string;
    dateRangeError?: string;
    isSaveEnabled: boolean;
    onClose: () => void;
    onStartDateChange: (date: Date | null) => void;
    onEndDateChange: (date: Date | null) => void;
    onStartAtDayStartChange: (checked: boolean) => void;
    onEndAtDayEndChange: (checked: boolean) => void;
    onSave: () => void;
}

export const EditDatesDialog: React.FC<EditDatesDialogProps> = ({
                                                                    isOpen,
                                                                    albumName,
                                                                    startDate,
                                                                    endDate,
                                                                    startAtDayStart,
                                                                    endAtDayEnd,
                                                                    isLoading,
                                                                    errorCode,
                                                                    dateRangeError,
                                                                    isSaveEnabled,
                                                                    onClose,
                                                                    onStartDateChange,
                                                                    onEndDateChange,
                                                                    onStartAtDayStartChange,
                                                                    onEndAtDayEndChange,
                                                                    onSave,
                                                                }) => {
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    // The dateError and dateHelperText should only show if there's a date range error,
    // not just because a date is null. Null dates disable the save button.
    const dateError = !!dateRangeError;
    const dateHelperText = dateRangeError || "";

    return (
        <Dialog
            open={isOpen}
            onClose={onClose}
            fullWidth
            fullScreen={isMobile}
            maxWidth="md"
        >
            <Box sx={{
                height: '4px',
                marginTop: '0px !important',
            }}>
                {isLoading && <LinearProgress sx={{
                    borderRadius: {
                        sm: '4px 4px 0px 0px'
                    },
                }}/>}
            </Box>
            <DialogTitle>Edit dates of {albumName}</DialogTitle>
            <IconButton
                aria-label="close"
                onClick={onClose}
                color='primary'
                sx={{
                    position: 'absolute',
                    right: 8,
                    top: 8,
                    color: (theme) => theme.palette.grey[500],
                }}
            >
                <Close/>
            </IconButton>
            <DialogContent>
                <Grid container spacing={2} alignItems='center'>
                    <Grid sm={12} xs={12}>
                        {errorCode && <Alert severity="error">
                            {errorCode}
                        </Alert>}
                    </Grid>
                    <DateRangePicker
                        startDate={startDate}
                        endDate={endDate}
                        startAtDayStart={startAtDayStart || false}
                        endAtDayEnd={endAtDayEnd || false}
                        onStartDateChange={onStartDateChange}
                        onEndDateChange={onEndDateChange}
                        onStartsAtStartOfTheDayChange={onStartAtDayStartChange}
                        onEndsAtEndOfTheDayChange={onEndAtDayEndChange}
                        disabled={isLoading}
                        dateError={dateError}
                        dateHelperText={dateHelperText}
                    />
                </Grid>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose} color="info" disabled={isLoading}>
                    Cancel
                </Button>
                <Button
                    onClick={onSave}
                    color="primary"
                    variant="contained"
                    disabled={!isSaveEnabled}
                >
                    Save
                </Button>
            </DialogActions>
        </Dialog>
    );
};
