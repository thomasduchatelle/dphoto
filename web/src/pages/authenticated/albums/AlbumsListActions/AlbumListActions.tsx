import {OwnerSelector, OwnerSelectorProps} from "./OwnerSelector";
import {Box, IconButton, Tooltip} from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import DeleteIcon from "@mui/icons-material/Delete";
import EditIcon from "@mui/icons-material/Edit";
import {CreateAlbumControls} from "../../../../core/catalog";
import {useSearchParams} from "react-router-dom";


export default function AlbumListActions({
                                             openDialogForCreateAlbum,
                                             openDeleteAlbumDialog,
                                             openEditDatesDialog,
                                             ...props
                                         }: OwnerSelectorProps & CreateAlbumControls & {
    openDeleteAlbumDialog: () => void
    openEditDatesDialog: () => void
}) {
    const [urlParams] = useSearchParams()
    const features = urlParams.get('features');
    const isEditAlbumFeatureEnabled = features && features.includes('edit-album');

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
            <Tooltip title={isEditAlbumFeatureEnabled ? "" : "Function not yet available, stay tuned !"}>
                <span>
                    <IconButton color="primary" size="large" onClick={openEditDatesDialog} disabled={!isEditAlbumFeatureEnabled}>
                        <EditIcon/>
                    </IconButton>
                </span>
            </Tooltip>
            <IconButton color="primary" size="large" onClick={openDeleteAlbumDialog}>
                <DeleteIcon/>
            </IconButton>
        </Box>
    )
}
