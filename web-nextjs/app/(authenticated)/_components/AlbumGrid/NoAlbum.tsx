'use client';

import CollectionsIcon from '@mui/icons-material/Collections';
import {Button} from '@mui/material';
import {PageMessage} from '@/components/PageMessage';

export interface NoAlbumProps {
    onCreateAlbum?: () => void;
}

export const NoAlbum = ({onCreateAlbum}: NoAlbumProps) => {
    return (
        <PageMessage
            icon={<CollectionsIcon/>}
            title="No Albums Found"
            message="Create your first album to get started organizing your photos."
        >
            {onCreateAlbum && (
                <Button variant="contained" onClick={onCreateAlbum}>
                    Create Album
                </Button>
            )}
        </PageMessage>
    );
};
