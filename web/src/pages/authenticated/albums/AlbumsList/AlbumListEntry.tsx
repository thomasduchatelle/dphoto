import {Album} from "../../../../core/catalog";
import {Avatar, Badge, ListItemAvatar, ListItemButton, ListItemSecondaryAction, ListItemText} from "@mui/material";
import {Link} from "react-router-dom";
import {toLocaleDateWithDay} from "../../../../core/utils/date-utils";
import {HotIndicator} from "./HotIndicator";

export function AlbumListEntry({album, selected}: {
    album: Album
    selected: boolean
}) {
    return <ListItemButton
        divider={false}
        selected={selected}
        to={`/albums/${album.albumId.owner}/${album.albumId.folderName}`}
        component={Link}
    >
        <ListItemAvatar>
            <Badge badgeContent={album.totalCount} color="info" max={999}>
                {(album.name.length > 2
                        && (
                            <Avatar>{`${album.name.charAt(0) + album.name.charAt(1)}`.toUpperCase()}</Avatar>))
                    || (<Avatar>??</Avatar>)}
            </Badge>
        </ListItemAvatar>
        <ListItemText
            primary={album.name}
            secondary={`${toLocaleDateWithDay(album.start)} → ${toLocaleDateWithDay(album.end)}`}
            sx={{
                pl: 1
            }}
        />
        <ListItemSecondaryAction title={`${+album.temperature.toFixed(2)} medias per day`}>
            <HotIndicator count={album.totalCount} relativeTemperature={album.relativeTemperature}/>
        </ListItemSecondaryAction>
    </ListItemButton>;
}