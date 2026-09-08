'use client';

import {Box} from '@mui/material';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import Link from '@/components/Link';
import {Album, AlbumId} from '@/domains/catalog/language';
import {albumUrl} from '@/domains/catalog/navigation/album-url';
import {AlbumCard} from '@/components/AlbumCard';

export interface NeighbourAlbumLinkProps {
    album: Album;
    onShare: (albumId: AlbumId) => void;
}

export function NeighbourAlbumLink({album, onShare}: NeighbourAlbumLinkProps) {
    return (
        <Box
            component={Link}
            href={albumUrl(album.albumId)}
            prefetch={false}
            sx={{
                display: 'flex',
                alignItems: 'center',
                gap: 1,
                textDecoration: 'none',
                maxWidth: 340,
            }}
        >
            <Box sx={{flex: 1, minWidth: 0}}>
                <AlbumCard album={album} compact onShare={onShare}/>
            </Box>
            <ChevronRightIcon sx={{color: 'rgba(255,255,255,0.35)', fontSize: 22, flexShrink: 0}}/>
        </Box>
    );
}
