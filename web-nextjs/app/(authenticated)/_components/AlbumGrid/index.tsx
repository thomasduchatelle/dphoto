'use client';

import {Box} from '@mui/material';
import Link from 'next/link';
import {Album, AlbumId} from '@/domains/catalog/language/catalog-state';
import {AlbumCard} from '../AlbumCard';

export interface AlbumGridProps {
    albums: Album[];
    onShare: (albumId: AlbumId) => void;
}

export const AlbumGrid = ({albums, onShare}: AlbumGridProps) => {
    return (
        <Box
            component="section"
            aria-label="Album list"
            sx={{
                display: 'grid',
                gridTemplateColumns: {
                    xs: '1fr',
                    sm: 'repeat(auto-fill, minmax(257px, 1fr))',
                },
                gap: 4,
                width: '100%',
                maxWidth: 1920,
                mx: 'auto',
                overflow: 'hidden',
            }}
        >
            {albums.map(album => (
                <Link
                    key={`${album.albumId.owner}-${album.albumId.folderName}`}
                    href={`/albums/${album.albumId.owner}/${album.albumId.folderName}`}
                    prefetch={false}
                    style={{textDecoration: 'none', display: 'block', minWidth: 0}}
                >
                    <AlbumCard album={album} onShare={onShare}/>
                </Link>
            ))}
        </Box>
    );
};
