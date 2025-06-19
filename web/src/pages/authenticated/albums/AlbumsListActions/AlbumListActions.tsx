import {OwnerSelector, OwnerSelectorProps} from "./OwnerSelector";
import {Box, IconButton} from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import DeleteIcon from "@mui/icons-material/Delete";
import {AlbumId, CreateAlbumControls} from "../../../../core/catalog";


export default function AlbumListActions({
                                             openDialogForCreateAlbum,
                                             openDeleteAlbumDialog,
                                             openEditAlbumDatesDialog,
                                             selectedAlbumId,
                                             ...props
                                         }: OwnerSelectorProps & CreateAlbumControls & {
    openDeleteAlbumDialog: () => void;
    openEditAlbumDatesDialog: (albumId: AlbumId) => void;
    selectedAlbumId?: AlbumId;
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
            >
                <DeleteIcon/>
            </IconButton>
            <IconButton color="primary" size="large" onClick={openDeleteAlbumDialog}>
                <DeleteIcon/>
            </IconButton>
        </Box>
    )
}
