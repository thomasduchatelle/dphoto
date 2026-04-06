'use client';

import AddPhotoAlternateIcon from '@mui/icons-material/AddPhotoAlternate';
import {Button} from '@mui/material';
import {PageMessage} from '@/components/PageMessage';
import Link from '@/components/Link';

export const NoMedia = () => {
    return (
        <PageMessage
            icon={<AddPhotoAlternateIcon/>}
            title="No Medias"
            message="Upload medias to this album to see them displayed here."
        >
            <Button
                component={Link}
                href="/"
                prefetch={false}
                variant="text"
            >
                Back to Albums
            </Button>
        </PageMessage>
    );
};
