'use client';

import {Box, Button, IconButton, Stack, Tooltip, Typography} from '@mui/material';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import Link from '@/components/Link';
import {Album} from '@/domains/catalog/language';
import PlayArrowIcon from "@mui/icons-material/PlayArrow";
import ShareIcon from "@mui/icons-material/Share";
import DriveFileRenameOutlineIcon from "@mui/icons-material/DriveFileRenameOutline";
import CalendarMonthIcon from "@mui/icons-material/CalendarMonth";
import DeleteIcon from "@mui/icons-material/Delete";
import AddIcon from "@mui/icons-material/Add";
import {toLocaleDateWithDay} from "@/libs/dates";

export interface AlbumHeaderProps {
    album: Album | undefined;
    isOwner: boolean;
    playHref: string | undefined;
    onShare: () => void;
    onRename: () => void;
    onEditDates: () => void;
    onDelete: () => void;
    onCreateAlbum: () => void;
}

export function AlbumHeader({album, isOwner, playHref, onShare, onRename, onEditDates, onDelete, onCreateAlbum}: AlbumHeaderProps) {
    return (
        <Box
            sx={(theme) => ({
                pr: {xs: 0, lg: theme.spacing(2)},
                py: {xs: 1.5, sm: 2.5},
                display: 'flex',
                alignItems: 'center',
                gap: {xs: 1, lg: 2},
            })}
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
                        {toLocaleDateWithDay(album.start)} – {toLocaleDateWithDay(album.end)} · {album.totalCount} photos
                    </Typography>
                )}
            </Box>
            <Stack direction="row" gap={1} alignItems="center" sx={{display: {xs: 'none', lg: 'flex'}, flexShrink: 0}}>
                <Button
                    variant="outlined"
                    color='secondary'
                    startIcon={<PlayArrowIcon/>}
                    component={Link}
                    href={playHref ?? ''}
                    prefetch={false}
                    disabled={!playHref}
                >
                    Play
                </Button>
                {isOwner && (
                    <>
                        <Button variant="outlined" color='secondary' startIcon={<AddIcon/>} onClick={onCreateAlbum}>Create</Button>
                        <Button variant="outlined" color='secondary' startIcon={<ShareIcon/>} onClick={onShare}>Share</Button>
                        <Button variant="outlined" color='secondary' startIcon={<DriveFileRenameOutlineIcon/>} onClick={onRename}>Rename</Button>
                        <Button variant="outlined" color='secondary' startIcon={<CalendarMonthIcon/>} onClick={onEditDates}>Dates</Button>
                        <Tooltip title="Delete album">
                            <IconButton color='secondary' onClick={onDelete} aria-label="Delete album">
                                <DeleteIcon/>
                            </IconButton>
                        </Tooltip>
                    </>
                )}
            </Stack>
        </Box>
    );
}
