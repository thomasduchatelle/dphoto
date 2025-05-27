import {OwnerSelector, OwnerSelectorProps} from "./OwnerSelector";
import {Box, IconButton, Tooltip} from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import SettingsIcon from "@mui/icons-material/Settings";
import {CreateAlbumControls} from "../../../../core/catalog/substates/CreateAlbumController";


export default function AlbumListActions({openDialogForCreateAlbum, ...props}: OwnerSelectorProps & CreateAlbumControls) {
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
            <IconButton color="primary" onClick={openDialogForCreateAlbum} size="large">
                <AddIcon/>
            </IconButton>
            <Tooltip title="Album management [Feature not yet available]">
                <span>
                    <IconButton color="primary" size="large" disabled={true}>
                        <SettingsIcon/>
                    </IconButton>
                </span>
            </Tooltip>
        </Box>
    )
}