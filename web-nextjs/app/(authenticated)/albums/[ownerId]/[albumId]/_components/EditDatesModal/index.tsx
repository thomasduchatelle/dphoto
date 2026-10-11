'use client';

import React from 'react';
import {
    Alert,
    Box,
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Grid,
    IconButton,
    LinearProgress,
    useMediaQuery,
    useTheme,
} from "@mui/material";
import {Close} from "@mui/icons-material";
import {LocalizationProvider} from "@mui/x-date-pickers";
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import {EditDatesDialogSelection} from "@/domains/catalog";
import {DateRangePicker} from "@/components/DateRangePicker";

export interface EditDatesDialogHandlers {
    onClose: () => void;
    onSave: () => Promise<void> | void;
    onStartDateChange: (date: Date | null) => void;
    onEndDateChange: (date: Date | null) => void;
    onStartAtDayStartChange: (checked: boolean) => void;
    onEndAtDayEndChange: (checked: boolean) => void;
}

export type EditDatesModalProps = EditDatesDialogSelection & EditDatesDialogHandlers;

export function EditDatesModal({
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
                                   onSave,
                                   onStartDateChange,
                                   onEndDateChange,
                                   onStartAtDayStartChange,
                                   onEndAtDayEndChange,
                               }: EditDatesModalProps) {
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    const handleSubmit = (event: React.FormEvent) => {
        event.preventDefault();
        if (isSaveEnabled) {
            void onSave();
        }
    };

    return (
        <LocalizationProvider dateAdapter={AdapterDayjs}>
            <Dialog
                open={isOpen}
                onClose={onClose}
                fullWidth
                fullScreen={isMobile}
                maxWidth="md"
                aria-labelledby="edit-dates-dialog-title"
            >
                <Box sx={{height: '4px', marginTop: '0px !important'}}>
                    {isLoading && (
                        <LinearProgress sx={{borderRadius: {sm: '4px 4px 0px 0px'}}}/>
                    )}
                </Box>
                <DialogTitle id="edit-dates-dialog-title">Edit dates of {albumName}</DialogTitle>
                <IconButton
                    aria-label="close"
                    onClick={onClose}
                    color="primary"
                    sx={{
                        position: 'absolute',
                        right: 8,
                        top: 8,
                        color: (theme) => theme.palette.grey[500],
                    }}
                >
                    <Close/>
                </IconButton>
                <form onSubmit={handleSubmit}>
                    <DialogContent>
                        <Grid container spacing={2} alignItems="center">
                            {errorCode && (
                                <Grid size={{sm: 12, xs: 12}}>
                                    <Alert severity="error">{errorCode}</Alert>
                                </Grid>
                            )}
                            <DateRangePicker
                                startDate={startDate}
                                endDate={endDate}
                                startAtDayStart={startAtDayStart}
                                endAtDayEnd={endAtDayEnd}
                                onStartDateChange={onStartDateChange}
                                onEndDateChange={onEndDateChange}
                                onStartsAtStartOfTheDayChange={onStartAtDayStartChange}
                                onEndsAtEndOfTheDayChange={onEndAtDayEndChange}
                                disabled={isLoading}
                                dateError={!!dateRangeError}
                                dateHelperText={dateRangeError}
                            />
                        </Grid>
                    </DialogContent>
                    <DialogActions>
                        <Button onClick={onClose} color="info" disabled={isLoading}>Cancel</Button>
                        <Button
                            type="submit"
                            color="primary"
                            variant="contained"
                            disabled={!isSaveEnabled}
                        >
                            Save
                        </Button>
                    </DialogActions>
                </form>
            </Dialog>
        </LocalizationProvider>
    );
}
