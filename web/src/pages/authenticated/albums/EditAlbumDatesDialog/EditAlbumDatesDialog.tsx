import React from "react";
import {Box, Button, Checkbox, Dialog, DialogActions, DialogContent, DialogTitle, FormControlLabel, TextField} from "@mui/material";
import {EditAlbumDatesDialogProps} from "./EditAlbumDatesDialogProps";

export const EditAlbumDatesDialog: React.FC<EditAlbumDatesDialogProps> = ({
                                                                              isOpen,
                                                                              albumName,
                                                                              startDate,
                                                                              endDate,
                                                                              isStartDateAtStartOfDay,
                                                                              isEndDateAtEndOfDay,
                                                                              onClose,
                                                                          }) => {
    return (
        <Dialog open={isOpen} onClose={onClose} maxWidth="sm" fullWidth>
            <DialogTitle>Edit Dates for "{albumName}"</DialogTitle>
            <DialogContent>
                <Box sx={{display: 'flex', flexDirection: 'column', gap: 2, mt: 1}}>
                    <TextField
                        label="Start Date"
                        type="date"
                        value={startDate}
                        InputLabelProps={{
                            shrink: true,
                        }}
                        fullWidth
                        disabled
                    />
                    <FormControlLabel
                        control={<Checkbox checked={isStartDateAtStartOfDay} disabled/>}
                        label="at the start of the day"
                    />
                    <TextField
                        label="End Date"
                        type="date"
                        value={endDate}
                        InputLabelProps={{
                            shrink: true,
                        }}
                        fullWidth
                        disabled
                    />
                    <FormControlLabel
                        control={<Checkbox checked={isEndDateAtEndOfDay} disabled/>}
                        label="at the end of the day"
                    />
                </Box>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose} color="info">
                    Cancel
                </Button>
            </DialogActions>
        </Dialog>
    );
};
