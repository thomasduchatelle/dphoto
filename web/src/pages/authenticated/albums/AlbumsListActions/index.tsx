import AddIcon from "@mui/icons-material/Add";
import SettingsIcon from '@mui/icons-material/Settings';
import {Box, IconButton, Tooltip} from "@mui/material";
import {OwnerSelector, OwnerSelectorProps} from "./OwnerSelector";

export default function AlbumListActions({...props}: OwnerSelectorProps) {
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
                <OwnerSelector {...props} />
            </Box>
            <Tooltip title="Feature not yet available">
                <span>
                    <IconButton color="primary" title="Create new album" onClick={doSomething} size="large" disabled={true}>
                        <AddIcon/>
                    </IconButton>
                </span>
            </Tooltip>
            <Tooltip title="Feature not yet available">
                <span>
                    <IconButton color="primary" title="Album management" onClick={doSomething} size="large" disabled={true}>
                        <SettingsIcon/>
                    </IconButton>
                </span>
            </Tooltip>
        </Box>
    )
}