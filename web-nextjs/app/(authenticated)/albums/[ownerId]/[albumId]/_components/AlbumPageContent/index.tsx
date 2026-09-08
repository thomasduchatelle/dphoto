'use client';

import {useEffect, useReducer, useRef, useState} from 'react';
import {notFound} from 'next/navigation';
import {Box, Button, IconButton, SpeedDial, SpeedDialAction, SpeedDialIcon, Stack, Typography} from '@mui/material';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import ShareIcon from '@mui/icons-material/Share';
import DriveFileRenameOutlineIcon from '@mui/icons-material/DriveFileRenameOutline';
import CalendarMonthIcon from '@mui/icons-material/CalendarMonth';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import Link from '@/components/Link';
import {catalogReducer, catalogThunks, CatalogViewerState} from '@/domains/catalog';
import {catalogViewerPageSelector} from '@/domains/catalog/navigation/selector-catalog-viewer-page';
import {useThunks} from '@/libs/dthunks/react';
import {ErrorMessage} from '@/components/ErrorMessage';
import {AlbumCard} from '../../../../../_components/AlbumCard';
import {AlbumMediaGrid, NeighbourAlbumLink} from '../AlbumMediaGrid';
import {Album} from '@/domains/catalog/language';

export interface AlbumPageContentProps {
    initialState: CatalogViewerState;
}

const fmt = (d: Date) => d.toLocaleDateString(undefined, {month: 'short', day: 'numeric', year: 'numeric'});

function albumHref(album: Album) {
    return `/albums/${album.albumId.owner}/${album.albumId.folderName}`;
}

function isSelected(album: Album, displayed: Album | undefined): boolean {
    return !!displayed
        && album.albumId.owner === displayed.albumId.owner
        && album.albumId.folderName === displayed.albumId.folderName;
}

export function AlbumPageContent({initialState}: AlbumPageContentProps) {
    const [state, dispatch] = useReducer(catalogReducer, initialState);

    const {onPageRefresh, loadAlbumPage, deleteAlbum, updateAlbumDates, submitCreateAlbum, saveAlbumName, grantAlbumAccess, revokeAlbumAccess, ...dispatchOnlyThunks} = catalogThunks;
    useThunks(dispatchOnlyThunks, {dispatch}, state);

    const {displayedAlbum, medias, mediasLoaded, albums, previousAlbum, nextAlbum, albumNotFound, error} = catalogViewerPageSelector(state);

    const [nextBannerVisible, setNextBannerVisible] = useState(false);
    const lastScrollY = useRef(0);
    const scrollRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const el = scrollRef.current;
        if (!el) return;
        const onScroll = () => {
            const y = el.scrollTop;
            setNextBannerVisible(y < lastScrollY.current && y > 80);
            lastScrollY.current = y;
        };
        el.addEventListener('scroll', onScroll, {passive: true});
        return () => el.removeEventListener('scroll', onScroll);
    }, []);

    if (error) {
        return <ErrorMessage error={error} title="Failed to load the album"/>;
    }

    if (albumNotFound) {
        notFound();
    }

    if (!mediasLoaded || medias.length === 0) {
        return (
            <NoMediaLayout
                albums={albums}
                displayedAlbum={displayedAlbum}
                previousAlbum={previousAlbum}
                nextAlbum={nextAlbum}
            />
        );
    }

    return (
        <Box
            ref={scrollRef}
            sx={{height: '100vh', overflowY: 'auto', display: 'flex', flexDirection: 'column'}}
        >
            {/* Next-album peek banner — mobile only (xs/sm/md), slides in on scroll-up */}
            {nextAlbum && (
                <NextAlbumBanner album={nextAlbum} visible={nextBannerVisible}/>
            )}

            <Box sx={{display: 'flex', flex: 1, mt: {xs: '56px', sm: '64px'}}}>

                {/* Left rail — lg+ only, scrolls with page */}
                <AlbumRail albums={albums} displayedAlbum={displayedAlbum}/>

                {/* Main column */}
                <Box sx={{flex: 1, minWidth: 0, display: 'flex', flexDirection: 'column'}}>

                    {/* Album info band */}
                    <Box
                        sx={{
                            px: {xs: 1.5, sm: 3, md: 5},
                            py: {xs: 1.5, sm: 2.5},
                            display: 'flex',
                            alignItems: 'center',
                            flexDirection: 'row',
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
                                {displayedAlbum?.name}
                            </Typography>
                            {displayedAlbum && (
                                <Typography variant="body1" sx={{fontSize: '0.85rem'}}>
                                    {fmt(displayedAlbum.start)} – {fmt(displayedAlbum.end)} · {displayedAlbum.totalCount} photos
                                </Typography>
                            )}
                        </Box>

                        {/* Desktop action buttons — lg+ */}
                        <Stack
                            direction="row"
                            gap={1}
                            sx={{display: {xs: 'none', lg: 'flex'}, flexShrink: 0}}
                        >
                            <ActionButtons/>
                        </Stack>
                    </Box>

                    {/* Media grid — next/prev navigation cards rendered inside on mobile */}
                    <Box sx={{px: {xs: 0, sm: 1, md: 3}, py: 3, flex: 1}}>
                        <AlbumMediaGrid
                            medias={medias}
                            nextAlbum={nextAlbum}
                            previousAlbum={previousAlbum}
                        />
                    </Box>
                </Box>
            </Box>

            {/* Mobile FAB — xs/sm/md only */}
            <SpeedDial
                ariaLabel="Album actions"
                sx={{
                    display: {xs: 'flex', lg: 'none'},
                    position: 'fixed',
                    bottom: 24,
                    right: 20,
                    '& .MuiSpeedDial-fab': {
                        bgcolor: '#185986',
                        '&:hover': {bgcolor: '#1d6fa3'},
                    },
                }}
                icon={<SpeedDialIcon openIcon={<PlayArrowIcon/>} icon={<PlayArrowIcon/>}/>}
            >
                <SpeedDialAction icon={<ShareIcon/>} tooltipTitle="Share" tooltipOpen/>
                <SpeedDialAction icon={<DriveFileRenameOutlineIcon/>} tooltipTitle="Rename" tooltipOpen/>
                <SpeedDialAction icon={<CalendarMonthIcon/>} tooltipTitle="Dates" tooltipOpen/>
            </SpeedDial>
        </Box>
    );
}

export function AlbumRail({albums, displayedAlbum}: {albums: Album[]; displayedAlbum: Album | undefined}) {
    return (
        <Box
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
                    href={albumHref(album)}
                    prefetch={false}
                    sx={{
                        display: 'block',
                        textDecoration: 'none',
                        outline: isSelected(album, displayedAlbum) ? '2px solid #4a9ece' : 'none',
                        borderRadius: 1,
                        overflow: 'hidden',
                    }}
                >
                    <AlbumCard album={album} onShare={() => {}} compact/>
                </Box>
            ))}
        </Box>
    );
}

function ActionButtons() {
    return (
        <>
            <Button variant="outlined" startIcon={<PlayArrowIcon/>} disabled sx={disabledButtonSx}>Play</Button>
            <Button variant="outlined" startIcon={<ShareIcon/>} disabled sx={disabledButtonSx}>Share</Button>
            <Button variant="outlined" startIcon={<DriveFileRenameOutlineIcon/>} disabled sx={disabledButtonSx}>Rename</Button>
            <Button variant="outlined" startIcon={<CalendarMonthIcon/>} disabled sx={disabledButtonSx}>Dates</Button>
        </>
    );
}

const disabledButtonSx = {
    borderColor: 'rgba(255,255,255,0.45)',
    color: 'white',
    '&.Mui-disabled': {borderColor: 'rgba(255,255,255,0.2)', color: 'rgba(255,255,255,0.3)'},
};

function NextAlbumBanner({album, visible}: {album: Album; visible: boolean}) {
    return (
        <Box
            component={Link}
            href={albumHref(album)}
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
                    <AlbumCard album={album} onShare={() => {}} compact/>
                </Box>
                <ChevronRightIcon sx={{color: 'rgba(255,255,255,0.4)', fontSize: 20, flexShrink: 0}}/>
            </Box>
        </Box>
    );
}

interface NoMediaLayoutProps {
    albums: Album[];
    displayedAlbum: Album | undefined;
    previousAlbum: Album | undefined;
    nextAlbum: Album | undefined;
}

function NoMediaLayout({albums, displayedAlbum, previousAlbum, nextAlbum}: NoMediaLayoutProps) {
    const hasNeighbours = !!(nextAlbum || previousAlbum);

    return (
        <Box sx={{display: 'flex', minHeight: '100vh', mt: {xs: '56px', sm: '64px'}}}>

            {/* Left rail — lg+ */}
            <AlbumRail albums={albums} displayedAlbum={displayedAlbum}/>

            {/* Centre column */}
            <Box
                sx={{
                    flex: 1,
                    display: 'flex',
                    flexDirection: 'column',
                    minHeight: 0,
                }}
            >
                {/* Album header — same as media page */}
                <Box
                    sx={{
                        px: {xs: 1.5, sm: 3, md: 5},
                        py: {xs: 1.5, sm: 2.5},
                        display: 'flex',
                        alignItems: 'center',
                        flexDirection: 'row',
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
                            {displayedAlbum?.name}
                        </Typography>
                        {displayedAlbum && (
                            <Typography variant="body1" sx={{fontSize: '0.85rem'}}>
                                {fmt(displayedAlbum.start)} – {fmt(displayedAlbum.end)}
                            </Typography>
                        )}
                    </Box>
                    <Stack direction="row" gap={1} sx={{display: {xs: 'none', lg: 'flex'}, flexShrink: 0}}>
                        <ActionButtons/>
                    </Stack>
                </Box>

                {/* Empty-album message — centred, main point of the page */}
                <Box
                    sx={{
                        flex: 1,
                        display: 'flex',
                        flexDirection: 'column',
                        alignItems: 'center',
                        justifyContent: hasNeighbours ? 'flex-start' : 'center',
                        px: {xs: 2, sm: 4},
                        py: 6,
                        gap: 4,
                        textAlign: 'center',
                    }}
                >
                    <Box>
                        <Typography variant="h1" sx={{mb: 1, opacity: 0.5}}>No photos yet</Typography>
                        <Typography variant="body1">Upload medias to this album to see them here.</Typography>
                    </Box>

                    {/* Neighbour album cards */}
                    {hasNeighbours && (
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
                                    <NeighbourAlbumLink album={nextAlbum}/>
                                </Box>
                            )}
                            {previousAlbum && (
                                <Box sx={{width: {xs: '100%', sm: 320}, flexShrink: 0}}>
                                    <NeighbourAlbumLink album={previousAlbum}/>
                                </Box>
                            )}
                        </Box>
                    )}
                </Box>
            </Box>

            {/* Mobile FAB — xs/sm/md only */}
            <SpeedDial
                ariaLabel="Album actions"
                sx={{
                    display: {xs: 'flex', lg: 'none'},
                    position: 'fixed',
                    bottom: 24,
                    right: 20,
                    '& .MuiSpeedDial-fab': {
                        bgcolor: '#185986',
                        '&:hover': {bgcolor: '#1d6fa3'},
                    },
                }}
                icon={<SpeedDialIcon openIcon={<PlayArrowIcon/>} icon={<PlayArrowIcon/>}/>}
            >
                <SpeedDialAction icon={<ShareIcon/>} tooltipTitle="Share" tooltipOpen/>
                <SpeedDialAction icon={<DriveFileRenameOutlineIcon/>} tooltipTitle="Rename" tooltipOpen/>
                <SpeedDialAction icon={<CalendarMonthIcon/>} tooltipTitle="Dates" tooltipOpen/>
            </SpeedDial>
        </Box>
    );
}
