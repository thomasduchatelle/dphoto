'use client';

import {Box, IconButton} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import Link from 'next/link';
import {Album, AlbumFilterCriterion, AlbumFilterEntry, AlbumId} from '@/domains/catalog/language/catalog-state';
import {albumUrl} from '@/domains/catalog/navigation/album-url';
import {AlbumCard} from '@/components/AlbumCard';
import {NoAlbum} from './NoAlbum';
import {AlbumFilterControl} from "../AlbumFilterControl";

export interface AlbumGridProps {
    albums: Album[];
    onShare: (albumId: AlbumId) => void;
    onCreateAlbum?: () => void;
    filterOptions: AlbumFilterEntry[];
    activeFilter: AlbumFilterEntry;
    onFilterChange: (filter: AlbumFilterCriterion) => void;
}

export const AlbumGrid = ({albums, onShare, onCreateAlbum, filterOptions, activeFilter, onFilterChange}: AlbumGridProps) => {
    if (albums.length === 0) {
        return <NoAlbum onCreateAlbum={onCreateAlbum}/>;
    }

    return (
        <Box>
            <Box sx={{
                width: '100%',
                maxWidth: 1920,
                mx: 'auto',
                marginBottom: 4,
                display: 'flex',
                alignItems: 'center',
                gap: 1,
            }}>
                <AlbumFilterControl
                    filterOptions={filterOptions}
                    activeFilter={activeFilter}
                    onFilterChange={onFilterChange}
                />
                {onCreateAlbum && (
                    <IconButton
                        color="primary"
                        size="large"
                        onClick={onCreateAlbum}
                        aria-label="Create album"
                    >
                        <AddIcon/>
                    </IconButton>
                )}
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
                        href={albumUrl(album.albumId)}
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
