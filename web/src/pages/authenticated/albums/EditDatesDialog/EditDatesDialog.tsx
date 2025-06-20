import React from "react";
import {Box, Button, CircularProgress, Dialog, DialogActions, DialogContent, DialogTitle, IconButton, useMediaQuery, useTheme} from "@mui/material";
import {Close} from "@mui/icons-material";
import Grid from "@mui/material/Unstable_Grid2";
import {DateRangePicker} from "../DateRangePicker";

interface EditDatesDialogProps {
    isOpen: boolean;
    albumName: string;
    startDate: Date;
    endDate: Date;
    startAtDayStart: boolean;
    endAtDayEnd: boolean;
    isLoading: boolean;
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
                                                                    onClose,
                                                                    onStartDateChange,
                                                                    onEndDateChange,
                                                                    onStartAtDayStartChange,
                                                                    onEndAtDayEndChange,
                                                                    onSave,
                                                                }) => {
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

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
            </Box>
            <DialogTitle>{albumName}</DialogTitle>
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
                    <DateRangePicker
                        startDate={startDate}
                        endDate={endDate}
                        startAtDayStart={startAtDayStart}
                        endAtDayEnd={endAtDayEnd}
                        onStartDateChange={onStartDateChange}
                        onEndDateChange={onEndDateChange}
                        onStartsAtStartOfTheDayChange={onStartAtDayStartChange}
                        onEndsAtEndOfTheDayChange={onEndAtDayEndChange}
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
                    disabled={isLoading}
                    startIcon={isLoading ? <CircularProgress size={20} /> : undefined}
                >
                    {isLoading ? 'Saving...' : 'Save'}
                </Button>
            </DialogActions>
        </Dialog>
    );
};
