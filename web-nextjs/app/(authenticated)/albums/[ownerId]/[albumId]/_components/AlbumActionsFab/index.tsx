'use client';

import {SpeedDial, SpeedDialAction, SpeedDialIcon} from '@mui/material';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import ShareIcon from '@mui/icons-material/Share';
import DriveFileRenameOutlineIcon from '@mui/icons-material/DriveFileRenameOutline';
import CalendarMonthIcon from '@mui/icons-material/CalendarMonth';
import DeleteIcon from '@mui/icons-material/Delete';

export interface AlbumActionsFabProps {
    isOwner: boolean;
    onRename: () => void;
    onEditDates: () => void;
    onDelete: () => void;
}

export function AlbumActionsFab({isOwner, onRename, onEditDates, onDelete}: AlbumActionsFabProps) {
    return (
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
            {isOwner && [
                <SpeedDialAction key="share" icon={<ShareIcon/>} tooltipTitle="Share" tooltipOpen FabProps={{disabled: true}}/>,
                <SpeedDialAction key="rename" icon={<DriveFileRenameOutlineIcon/>} tooltipTitle="Rename" tooltipOpen onClick={onRename}/>,
                <SpeedDialAction key="dates" icon={<CalendarMonthIcon/>} tooltipTitle="Dates" tooltipOpen onClick={onEditDates}/>,
                <SpeedDialAction key="delete" icon={<DeleteIcon/>} tooltipTitle="Delete" tooltipOpen onClick={onDelete}/>,
            ]}
        </SpeedDial>
    );
}
