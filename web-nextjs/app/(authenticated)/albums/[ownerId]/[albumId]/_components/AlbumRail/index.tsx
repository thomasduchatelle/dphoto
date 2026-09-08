'use client';

import {Box, Typography} from '@mui/material';
import Link from '@/components/Link';
import {Album, AlbumId, albumIdEquals} from '@/domains/catalog/language';
import {albumUrl} from '@/domains/catalog/navigation/album-url';
import {AlbumCard} from '@/components/AlbumCard';

export interface AlbumRailProps {
    albums: Album[];
    displayedAlbumId: AlbumId | undefined;
}

export function AlbumRail({albums, displayedAlbumId}: AlbumRailProps) {
    return (
        <Box
            component="nav"
            aria-label="Albums"
            sx={{
                display: {xs: 'none', lg: 'flex'},
                flexDirection: 'column',
                width: 450,
                flexShrink: 0,
                borderRight: '1px solid rgba(255,255,255,0.07)',
                bgcolor: 'rgba(0,10,20,0.5)',
                p: 2,
                gap: 1,
            }}
        >
            <Typography
                sx={{
                    fontSize: '0.65rem',
                    letterSpacing: '0.14em',
                    textTransform: 'uppercase',
                    color: 'rgba(255,255,255,0.35)',
                    mb: 0.5,
                }}
            >
                Albums
            </Typography>
            {albums.map((album) => (
                <Box
                    key={`${album.albumId.owner}/${album.albumId.folderName}`}
                    component={Link}
                    href={albumUrl(album.albumId)}
                    prefetch={false}
                    sx={{
                        display: 'block',
                        textDecoration: 'none',
                        outline: albumIdEquals(album.albumId, displayedAlbumId) ? '2px solid #4a9ece' : 'none',
                        borderRadius: 1,
                        overflow: 'hidden',
                    }}
                >
                    <AlbumCard album={album} compact/>
                </Box>
            ))}
        </Box>
    );
}
