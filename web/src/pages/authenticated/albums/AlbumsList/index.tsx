import {
  Avatar,
  Badge,
  Box,
  List,
  ListItem,
  ListItemAvatar,
  ListItemButton,
  ListItemSecondaryAction,
  ListItemText,
  Skeleton
} from "@mui/material";
import {Link} from "react-router-dom";
import {toLocaleDateWithDay} from "../../../../core/date-utils";
import {Album, AlbumId} from "../logic";
import {HotIndicator} from "./HotIndicator";

const AlbumsList = ({albums, loaded, selected, onSelect}: {
  albums: Album[]
  loaded: boolean
  selected?: Album
  onSelect(albumId: AlbumId): void
}) => {
  return (
    <Box sx={{overflow: 'auto'}}>
      <List>
        {!loaded && (Array.from(Array(4).keys()).map((v, index) => (
          <ListItem
            key={index}
            divider={index < 2}>
            <ListItemAvatar>
              <Skeleton variant="circular" width={35} height={35} animation='wave'/>
            </ListItemAvatar>
            <ListItemText
              primary={<Skeleton variant="text" width={200} animation='wave'/>}
              secondary={<Skeleton variant="text" width={300} animation='wave'/>}
            />
            <ListItemSecondaryAction>
              <Skeleton variant="circular" width={35} height={35} animation='wave'/>
            </ListItemSecondaryAction>
          </ListItem>
        )))}

        {albums.map((album, index) => (
          <ListItemButton
            divider={index + 1 < albums.length}
            selected={album.albumId.owner === selected?.albumId?.owner && album.albumId.folderName === selected?.albumId?.folderName}
            to={`/albums/${album.albumId.owner}/${album.albumId.folderName}`}
            component={Link}
            key={album.albumId.owner + album.albumId.folderName}
          >
            <ListItemAvatar>
              <Badge badgeContent={album.totalCount} color="info" max={999}>
                {(album.name.length > 2
                    && (<Avatar>{`${album.name.charAt(0) + album.name.charAt(1)}`.toUpperCase()}</Avatar>))
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
          </ListItemButton>
        ))}
      </List>
    </Box>
  )
}

export default AlbumsList
