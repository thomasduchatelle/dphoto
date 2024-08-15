import AddIcon from "@mui/icons-material/Add";
import SettingsIcon from '@mui/icons-material/Settings';
import {Box, IconButton} from "@mui/material";
import AlarmUserSelector from "./AlbumUserSelector";

export default function AlbumListActions() {
    const doSomething = () => {
        console.log("creating a new album")
    }

    return (
        <Box sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            '& > :not(style)': {mt: 1, mb: 1},
        }}>
            <Box sx={{mr: 2}}>
                <AlarmUserSelector/>
            </Box>
            <IconButton color="primary" title="Create new album" onClick={doSomething} size="large">
                <AddIcon/>
            </IconButton>
            <IconButton color="primary" title="Album management" onClick={doSomething} size="large">
                <SettingsIcon/>
            </IconButton>
        </Box>
    )
}