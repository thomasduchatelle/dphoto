'use client';

import {Box} from '@mui/material';
import Link from 'next/link';
import {Album, AlbumFilterEntry, AlbumId} from '@/domains/catalog/language/catalog-state';
import {AlbumCard} from '../AlbumCard';
import {NoAlbum} from './NoAlbum';
import {AlbumFilterControl} from "../AlbumFilterControl";

export interface AlbumGridProps {
    albums: Album[];
    onShare: (albumId: AlbumId) => void;
    onCreateAlbum?: () => void;
    filterOptions: AlbumFilterEntry[];
    activeFilter: AlbumFilterEntry;
    onFilterChange: (filter: AlbumFilterEntry) => void;
}

export const AlbumGrid = ({albums, onShare, onCreateAlbum, filterOptions, activeFilter, onFilterChange}: AlbumGridProps) => {
    if (albums.length === 0) {
        return <NoAlbum onCreateAlbum={onCreateAlbum}/>;
    }

    return (
        <Box>
            <Box sx={{width: '100%', maxWidth: 1920, mx: 'auto', marginBottom: 4}}>
                <AlbumFilterControl
                    filterOptions={filterOptions}
                    activeFilter={activeFilter}
                    onFilterChange={onFilterChange}
                />
            </Box>
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
        </Box>
    );
};
