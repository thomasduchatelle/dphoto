'use client';

import {Button, Dialog, DialogActions, DialogContent, DialogTitle, Typography} from '@mui/material';

export interface ShareAlbumModalStubProps {
    open: boolean;
    onClose: () => void;
}

export function ShareAlbumModalStub({open, onClose}: ShareAlbumModalStubProps) {
    return (
        <Dialog open={open} onClose={onClose}>
            <DialogTitle>Share — not implemented yet</DialogTitle>
            <DialogContent>
                <Typography>This feature is not implemented yet.</Typography>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Close</Button>
            </DialogActions>
        </Dialog>
    );
}
