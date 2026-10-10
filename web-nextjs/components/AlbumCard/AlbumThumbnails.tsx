'use client';

import {Box} from '@mui/material';
import {Album} from '@/domains/catalog/language/catalog-state';

export interface AlbumThumbnailsProps {
    album: Album;
    compact?: boolean;
}

export const AlbumThumbnails = ({album, compact = false}: AlbumThumbnailsProps) => (
    <Box sx={{display: 'grid', gridTemplateColumns: compact ? 'repeat(4, 1fr)' : {xs: 'repeat(4, 1fr)', sm: 'repeat(2, 1fr)'}, gap: 1}}>
        {[0, 1, 2, 3].map((i) => {
            const cover = album.covers[i];
            return (
                <Box
                    key={i}
                    className="album-photo"
                    sx={{
                        aspectRatio: '1',
                        background: 'linear-gradient(135deg, #1e3a5f, #2a4a6f)',
                        transition: 'filter 0.3s ease',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        fontSize: !cover ? 32 : 12,
                        opacity: !cover ? 0.15 : 1,
                        overflow: 'hidden',
                        position: 'relative',
                    }}
                >
                    {cover ? (
                        <img
                            src={cover.contentPath}
                            alt={`${album.name} cover ${i + 1}`}
                            style={{
                                width: '100%',
                                height: '100%',
                                objectFit: 'cover',
                            }}
                        />
                    ) : (
                        '🖼️'
                    )}
                </Box>
            );
        })}
    </Box>
);
