'use client';

import {Button, Dialog, DialogActions, DialogContent, DialogTitle, Typography} from '@mui/material';

export interface EditDatesModalStubProps {
    open: boolean;
    onClose: () => void;
}

export function EditDatesModalStub({open, onClose}: EditDatesModalStubProps) {
    return (
        <Dialog open={open} onClose={onClose}>
            <DialogTitle>Edit dates — not implemented yet</DialogTitle>
            <DialogContent>
                <Typography>This feature is not implemented yet.</Typography>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Close</Button>
            </DialogActions>
        </Dialog>
    );
}
