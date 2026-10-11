'use client';

import {Button, Dialog, DialogActions, DialogContent, DialogTitle, Typography} from '@mui/material';

export interface RenameAlbumModalStubProps {
    open: boolean;
    onClose: () => void;
}

export function RenameAlbumModalStub({open, onClose}: RenameAlbumModalStubProps) {
    return (
        <Dialog open={open} onClose={onClose}>
            <DialogTitle>Rename — not implemented yet</DialogTitle>
            <DialogContent>
                <Typography>This feature is not implemented yet.</Typography>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Close</Button>
            </DialogActions>
        </Dialog>
    );
}
