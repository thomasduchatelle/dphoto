import {Album} from "../../../../core/catalog";
import {
    Avatar,
    AvatarGroup,
    Badge,
    IconButton,
    ListItemAvatar,
    ListItemButton,
    ListItemSecondaryAction,
    ListItemText,
    Tooltip
} from "@mui/material";
import {Link} from "react-router-dom";
import {toLocaleDateWithDay} from "../../../../core/utils/date-utils";
import {HotIndicator} from "./HotIndicator";
import {Share} from "@mui/icons-material";
import React from "react";

export function AlbumListEntry({album, selected, onClickOnSharedWith}: {
    album: Album
    selected: boolean
    onClickOnSharedWith: (album: Album) => void
}) {

    const handleClickOnSharedWith = (evt: React.MouseEvent<HTMLElement>) => {
        evt.preventDefault()
        onClickOnSharedWith(album)
    }

    return <ListItemButton
        divider={false}
        selected={selected}
        to={`/albums/${album.albumId.owner}/${album.albumId.folderName}`}
        component={Link}
        sx={{
            borderRadius: '20px',
        }}
    >
        <ListItemAvatar>
            <Badge badgeContent={album.totalCount} color="info" max={999}
                   title={`${+album.temperature.toFixed(2)} medias per day`}>
                <HotIndicator count={album.totalCount} relativeTemperature={album.relativeTemperature}/>
            </Badge>
        </ListItemAvatar>
        <ListItemText
            primary={album.name}
            secondary={`${toLocaleDateWithDay(album.start)} → ${toLocaleDateWithDay(album.end)}`}
            sx={{
                pl: 1
            }}
        />
        {album.ownedBy ? (
            <Tooltip title={`Shared by ${album.ownedBy.name ?? "a friend"}`}>
                <AvatarGroup max={2} spacing='small' sx={{
                    '& .MuiAvatarGroup-avatar': {width: 32, height: 32, fontSize: "0.8em"},
                }}>
                    {album.ownedBy.users.map(user => (
                        <Avatar key={user.email} src={user.picture} alt={user.name}/>
                    ))}
                </AvatarGroup>
            </Tooltip>
        ) : (
            <ListItemSecondaryAction>
                <Badge
                    badgeContent={album.sharedWith.length ?? ''}
                    color="secondary"
                    max={9}
                    anchorOrigin={{"vertical": "bottom", "horizontal": "right"}}
                >
                    <IconButton onClick={handleClickOnSharedWith}>
                        <Share/>
                    </IconButton>
                </Badge>
            </ListItemSecondaryAction>
        )}
    </ListItemButton>;
}