import {OwnerSelector, OwnerSelectorProps} from "./OwnerSelector";
import {Box, IconButton, Menu, MenuItem} from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import DeleteIcon from "@mui/icons-material/Delete";
import EditIcon from "@mui/icons-material/Edit";
import {useState} from "react";

export default function AlbumListActions({
                                             openCreateDialog,
                                             openDeleteAlbumDialog,
                                             openEditDatesDialog,
                                             openEditNameDialog,
                                             displayedAlbumIdIsOwned = true,
                                             ...props
                                         }: OwnerSelectorProps & {
    openCreateDialog: () => void
    openDeleteAlbumDialog: () => void
    openEditDatesDialog: () => void
    openEditNameDialog: () => void
    displayedAlbumIdIsOwned: boolean
}) {
    const [editMenuAnchorEl, setEditMenuAnchorEl] = useState<null | HTMLElement>(null);
    const editMenuOpen = Boolean(editMenuAnchorEl);

    const handleEditClick = (event: React.MouseEvent<HTMLElement>) => {
        setEditMenuAnchorEl(event.currentTarget);
    };

    const handleEditMenuClose = () => {
        setEditMenuAnchorEl(null);
    };

    const handleEditDatesClick = () => {
        handleEditMenuClose();
        openEditDatesDialog();
    };

    const handleEditNameClick = () => {
        handleEditMenuClose();
        openEditNameDialog();
    };
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
            <IconButton color="primary" size="large" onClick={handleEditClick} disabled={!displayedAlbumIdIsOwned}>
                <EditIcon/>
            </IconButton>
            <Menu
                anchorEl={editMenuAnchorEl}
                open={editMenuOpen}
                onClose={handleEditMenuClose}
                MenuListProps={{
                    'aria-labelledby': 'edit-button',
                }}
            >
                <MenuItem onClick={handleEditDatesClick}>Edit Dates</MenuItem>
                <MenuItem onClick={handleEditNameClick} disabled>Edit Name</MenuItem>
            </Menu>
            <IconButton color="primary" size="large" onClick={openDeleteAlbumDialog}>
                <DeleteIcon/>
            </IconButton>
        </Box>
    )
}
