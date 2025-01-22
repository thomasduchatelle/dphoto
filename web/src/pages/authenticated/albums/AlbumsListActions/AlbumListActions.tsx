import {OwnerSelector, OwnerSelectorProps} from "./OwnerSelector";
import {Box, IconButton, Tooltip} from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import SettingsIcon from "@mui/icons-material/Settings";
import {useState} from "react";
import CreateAlbumDialog from "../CreateAlbumDialog";
import {OnCreateNewAlbumRequestType} from "../../../../core/catalog-react";

export interface AlbumListActionsCallBacks {
    onCreateNewAlbumRequest: OnCreateNewAlbumRequestType
}

export default function AlbumListActions({onCreateNewAlbumRequest, ...props}: OwnerSelectorProps & AlbumListActionsCallBacks) {
    const [createModal, setCreateModal] = useState(false)

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
            <IconButton color="primary" onClick={() => setCreateModal(true)} size="large">
                <AddIcon/>
            </IconButton>
            <Tooltip title="Album management [Feature not yet available]">
                <span>
                    <IconButton color="primary" size="large" disabled={true}>
                        <SettingsIcon/>
                    </IconButton>
                </span>
            </Tooltip>
            <CreateAlbumDialog
                open={createModal}
                onClose={() => setCreateModal(false)}
                onSubmit={(request) => {
                    setCreateModal(false)
                    return onCreateNewAlbumRequest(request)
                }}
            />
        </Box>
    )
}