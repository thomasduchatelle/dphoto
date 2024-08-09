import CollectionsIcon from "@mui/icons-material/Collections";
import AddIcon from "@mui/icons-material/Add";

import {Box, Fab} from "@mui/material";
import AlarmUserSelector from "./AlbumUserSelector";

export default function AlbumListActions() {
    return (
        <Box sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            '& > :not(style)': { m: 1, mr: 3 },
        }}>
            <AlarmUserSelector />
            <Fab variant="extended">
                <CollectionsIcon sx={{mr: 1}} />
                All Albums
            </Fab>
            <Fab color="secondary" aria-label="add">
                <AddIcon/>
            </Fab>
        </Box>
    )
}