'use client';

import {useEffect, useRef, useState} from 'react';
import {Box} from '@mui/material';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import Link from '@/components/Link';
import {Album} from '@/domains/catalog/language';
import {albumUrl} from '@/domains/catalog/navigation/album-url';
import {AlbumCard} from '@/components/AlbumCard';

export interface NextAlbumBannerProps {
    album: Album | undefined;
}

export function NextAlbumBanner({album}: NextAlbumBannerProps) {
    const [visible, setVisible] = useState(false);
    const lastScrollY = useRef(0);

    useEffect(() => {
        const onScroll = () => {
            setVisible(window.scrollY < lastScrollY.current && window.scrollY > 80);
            lastScrollY.current = window.scrollY;
        };
        window.addEventListener('scroll', onScroll, {passive: true});
        return () => window.removeEventListener('scroll', onScroll);
    }, []);

    if (!album) {
        return null;
    }

    return (
        <Box
            component={Link}
            href={albumUrl(album.albumId)}
            prefetch={false}
            sx={{
                display: {xs: 'flex', lg: 'none'},
                position: 'fixed',
                top: 0,
                left: 0,
                right: 0,
                zIndex: 1050,
                transform: visible ? 'translateY(0)' : 'translateY(-100%)',
                transition: 'transform 0.25s ease',
                flexDirection: 'column',
                textDecoration: 'none',
            }}
        >
            <Box sx={{height: {xs: '56px', sm: '64px'}, flexShrink: 0}}/>
            <Box
                sx={{
                    bgcolor: 'rgba(5,12,22,0.97)',
                    backdropFilter: 'blur(10px)',
                    borderBottom: '1px solid rgba(74,158,206,0.25)',
                    px: 2,
                    py: 1,
                    display: 'flex',
                    alignItems: 'center',
                    gap: 1.5,
                }}
            >
                <Box sx={{flex: 1, minWidth: 0, maxWidth: 400}}>
                    <AlbumCard album={album} compact/>
                </Box>
                <ChevronRightIcon sx={{color: 'rgba(255,255,255,0.4)', fontSize: 20, flexShrink: 0}}/>
            </Box>
        </Box>
    );
}
