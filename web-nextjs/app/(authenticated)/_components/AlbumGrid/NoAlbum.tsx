'use client';

import CollectionsIcon from '@mui/icons-material/Collections';
import {Button} from '@mui/material';
import {EmptyState, emptyStateButtonStyles} from '@/components/shared/EmptyState';

export interface NoAlbumProps {
    onCreateAlbum?: () => void;
}

export const NoAlbum = ({onCreateAlbum}: NoAlbumProps) => {
    return (
        <EmptyState
            icon={<CollectionsIcon/>}
            title="No Albums Found"
            message="Create your first album to get started organizing your photos."
        >
            {onCreateAlbum && (
                <Button
                    variant="contained"
                    onClick={onCreateAlbum}
                    sx={emptyStateButtonStyles.contained}
                >
                    Create Album
                </Button>
            )}
        </EmptyState>
    );
};
