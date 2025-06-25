import {OwnerSelector, OwnerSelectorProps} from "./OwnerSelector";
import {Box, IconButton} from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import DeleteIcon from "@mui/icons-material/Delete";
import EditIcon from "@mui/icons-material/Edit";
export default function AlbumListActions({
                                             openCreateDialog,
                                             openDeleteAlbumDialog,
                                             openEditDatesDialog,
                                             displayedAlbumIdIsOwned = true,
                                             ...props
                                         }: OwnerSelectorProps & {
    openCreateDialog: () => void
    openDeleteAlbumDialog: () => void
    openEditDatesDialog: () => void
    displayedAlbumIdIsOwned: boolean
}) {
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
            <IconButton color="primary" onClick={openCreateDialog} size="large">
                <AddIcon/>
            </IconButton>
            <IconButton color="primary" size="large" onClick={openEditDatesDialog} disabled={!displayedAlbumIdIsOwned}>
                <EditIcon/>
            </IconButton>
            <IconButton color="primary" size="large" onClick={openDeleteAlbumDialog}>
                <DeleteIcon/>
            </IconButton>
        </Box>
    )
}
