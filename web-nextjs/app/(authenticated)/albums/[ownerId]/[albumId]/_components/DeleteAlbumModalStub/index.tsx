'use client';

import {Button, Dialog, DialogActions, DialogContent, DialogTitle, Typography} from '@mui/material';

export interface DeleteAlbumModalStubProps {
    open: boolean;
    onClose: () => void;
}

export function DeleteAlbumModalStub({open, onClose}: DeleteAlbumModalStubProps) {
    return (
        <Dialog open={open} onClose={onClose}>
            <DialogTitle>Delete — not implemented yet</DialogTitle>
            <DialogContent>
                <Typography>This feature is not implemented yet.</Typography>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Close</Button>
            </DialogActions>
        </Dialog>
    );
}
