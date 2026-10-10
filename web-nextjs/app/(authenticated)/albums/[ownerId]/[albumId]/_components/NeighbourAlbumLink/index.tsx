'use client';

import {Box} from '@mui/material';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import Link from '@/components/Link';
import {Album} from '@/domains/catalog/language';
import {albumUrl} from '@/domains/catalog/navigation/album-url';
import {AlbumCard} from '@/components/AlbumCard';

export interface NeighbourAlbumLinkProps {
    album: Album;
    direction: 'previous' | 'next';
}

export function NeighbourAlbumLink({album, direction}: NeighbourAlbumLinkProps) {
    const chevronSx = {color: 'rgba(255,255,255,0.35)', fontSize: 22, flexShrink: 0};

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
                maxWidth: {sm: 340},
            }}
        >
            {direction === 'previous' && <ChevronLeftIcon sx={chevronSx}/>}
            <Box sx={{flex: 1, minWidth: 0}}>
                <AlbumCard album={album} compact/>
            </Box>
            {direction === 'next' && <ChevronRightIcon sx={chevronSx}/>}
        </Box>
    );
}
