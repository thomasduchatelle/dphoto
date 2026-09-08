'use client';

import AddPhotoAlternateIcon from '@mui/icons-material/AddPhotoAlternate';
import {Box} from '@mui/material';
import {PageMessage} from '@/components/PageMessage';
import {Album, AlbumId} from '@/domains/catalog/language';
import {NeighbourAlbumLink} from '../NeighbourAlbumLink';

export interface NoMediaProps {
    nextAlbum?: Album;
    previousAlbum?: Album;
    onShare: (albumId: AlbumId) => void;
}

export const NoMedia = ({nextAlbum, previousAlbum, onShare}: NoMediaProps) => {
    return (
        <PageMessage
            icon={<AddPhotoAlternateIcon/>}
            title="No Medias"
            message="Upload medias to this album to see them displayed here."
        >
            {(nextAlbum || previousAlbum) && (
                <Box
                    sx={{
                        display: 'flex',
                        flexDirection: {xs: 'column', sm: 'row'},
                        gap: 2,
                        width: '100%',
                        maxWidth: 700,
                        justifyContent: 'center',
                    }}
                >
                    {nextAlbum && (
                        <Box sx={{width: {xs: '100%', sm: 320}, flexShrink: 0}}>
                            <NeighbourAlbumLink album={nextAlbum} onShare={onShare}/>
                        </Box>
                    )}
                    {previousAlbum && (
                        <Box sx={{width: {xs: '100%', sm: 320}, flexShrink: 0}}>
                            <NeighbourAlbumLink album={previousAlbum} onShare={onShare}/>
                        </Box>
                    )}
                </Box>
            )}
        </PageMessage>
    );
};
