'use client';

import {Box, IconButton, Stack, Typography} from '@mui/material';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import Link from '@/components/Link';
import {Album} from '@/domains/catalog/language';
import {ActionButtons} from './ActionButtons';

export interface AlbumHeaderProps {
    album: Album | undefined;
}

const fmt = (d: Date) => d.toLocaleDateString(undefined, {month: 'short', day: 'numeric', year: 'numeric'});

export function AlbumHeader({album}: AlbumHeaderProps) {
    return (
        <Box
            sx={{
                px: {xs: 1.5, sm: 3, md: 5},
                py: {xs: 1.5, sm: 2.5},
                display: 'flex',
                alignItems: 'center',
                gap: {xs: 1, sm: 2},
                borderBottom: '1px solid rgba(255,255,255,0.07)',
            }}
        >
            <IconButton
                component={Link}
                href="/"
                prefetch={false}
                size="small"
                sx={{color: 'rgba(255,255,255,0.55)', flexShrink: 0}}
            >
                <ArrowBackIcon fontSize="small"/>
            </IconButton>
            <Box sx={{flex: 1, minWidth: 0}}>
                <Typography variant="h1" sx={{mb: 0.25}} noWrap>
                    {album?.name}
                </Typography>
                {album && (
                    <Typography variant="body1" sx={{fontSize: '0.85rem'}}>
                        {fmt(album.start)} – {fmt(album.end)} · {album.totalCount} photos
                    </Typography>
                )}
            </Box>
            <Stack direction="row" gap={1} sx={{display: {xs: 'none', lg: 'flex'}, flexShrink: 0}}>
                <ActionButtons/>
            </Stack>
        </Box>
    );
}
