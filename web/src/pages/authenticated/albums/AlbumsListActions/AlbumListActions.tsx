import {OwnerSelector, OwnerSelectorProps} from "./OwnerSelector";
import {Box, IconButton} from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import DeleteIcon from "@mui/icons-material/Delete";
import EditCalendarIcon from '@mui/icons-material/EditCalendar'; // New import
import {AlbumId, CreateAlbumControls} from "../../../../core/catalog";


export default function AlbumListActions({
                                             openDialogForCreateAlbum,
                                             openDeleteAlbumDialog,
                                             openEditAlbumDatesDialog, // New prop
                                             ...props
                                         }: OwnerSelectorProps & CreateAlbumControls & {
    openDeleteAlbumDialog: () => void;
    openEditAlbumDatesDialog: (albumId: AlbumId) => void; // New prop type
}) {
    // Determine if the "Edit Dates" button should be enabled.
    // It should be enabled only if an album is selected and it's owned by the current user.
    // For this story, we assume the selected album is always owned by the current user
    // if it's selected. The actual check for ownership will come in a later story.
    const isEditDatesEnabled = props.selected.criterion.selfOwned && props.selected.criterion.owners.length === 0;
    const selectedAlbumId = props.selected.criterion.selfOwned ? {owner: "myself", folderName: "current"} : undefined; // Placeholder for actual selected album ID

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
            <IconButton
                color="primary"
                size="large"
                onClick={() => {
                    if (selectedAlbumId) {
                        openEditAlbumDatesDialog(selectedAlbumId);
                    }
                }}
                disabled={!isEditDatesEnabled} // Disable if no album is selected or not self-owned
            >
                <EditCalendarIcon/>
            </IconButton>
            <IconButton color="primary" size="large" onClick={openDeleteAlbumDialog}>
                <DeleteIcon/>
            </IconButton>
        </Box>
    )
}
