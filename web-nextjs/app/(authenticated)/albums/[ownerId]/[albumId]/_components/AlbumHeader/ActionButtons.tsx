import {Button} from '@mui/material';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import ShareIcon from '@mui/icons-material/Share';
import DriveFileRenameOutlineIcon from '@mui/icons-material/DriveFileRenameOutline';
import CalendarMonthIcon from '@mui/icons-material/CalendarMonth';

const disabledButtonSx = {
    borderColor: 'rgba(255,255,255,0.45)',
    color: 'white',
    '&.Mui-disabled': {borderColor: 'rgba(255,255,255,0.2)', color: 'rgba(255,255,255,0.3)'},
};

export function ActionButtons() {
    return (
        <>
            <Button variant="outlined" startIcon={<PlayArrowIcon/>} disabled sx={disabledButtonSx}>Play</Button>
            <Button variant="outlined" startIcon={<ShareIcon/>} disabled sx={disabledButtonSx}>Share</Button>
            <Button variant="outlined" startIcon={<DriveFileRenameOutlineIcon/>} disabled sx={disabledButtonSx}>Rename</Button>
            <Button variant="outlined" startIcon={<CalendarMonthIcon/>} disabled sx={disabledButtonSx}>Dates</Button>
        </>
    );
}
