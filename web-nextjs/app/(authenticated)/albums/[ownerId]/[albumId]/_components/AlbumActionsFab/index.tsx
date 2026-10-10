'use client';

import {SpeedDial, SpeedDialAction, SpeedDialIcon} from '@mui/material';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import ShareIcon from '@mui/icons-material/Share';
import DriveFileRenameOutlineIcon from '@mui/icons-material/DriveFileRenameOutline';
import CalendarMonthIcon from '@mui/icons-material/CalendarMonth';

export function AlbumActionsFab() {
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
            <SpeedDialAction icon={<ShareIcon/>} tooltipTitle="Share" tooltipOpen/>
            <SpeedDialAction icon={<DriveFileRenameOutlineIcon/>} tooltipTitle="Rename" tooltipOpen/>
            <SpeedDialAction icon={<CalendarMonthIcon/>} tooltipTitle="Dates" tooltipOpen/>
        </SpeedDial>
    );
}
