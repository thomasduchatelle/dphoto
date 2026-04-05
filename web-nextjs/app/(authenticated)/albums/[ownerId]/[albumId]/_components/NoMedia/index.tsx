'use client';
import AddPhotoAlternateIcon from '@mui/icons-material/AddPhotoAlternate';
import {Button} from '@mui/material';
import {EmptyState, emptyStateButtonStyles} from '@/components/shared/EmptyState';
import Link from '@/components/Link';

export const NoMedia = () => {
    return (
        <EmptyState
            icon={<AddPhotoAlternateIcon/>}
            title="No Medias"
            message="Upload medias to this album to see them displayed here."
        >
            <Button
                component={Link}
                href="/"
                prefetch={false}
                variant="outlined"
                sx={emptyStateButtonStyles.outlined}
            >
                Back to Albums
            </Button>
        </EmptyState>
    );
};
