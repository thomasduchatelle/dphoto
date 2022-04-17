import {Alert, ImageList, ImageListItem, Skeleton} from "@mui/material";
import {Media} from "../../../../core/domain/catalog";
import useWindowDimensions from "../../../../core/window-utils";
import {ImageInList} from "./media-in-list";


const drawerWidth = 450;
const xsBreakpoint = 600;

export function MediaListComponent({medias, loaded, albumNotFound}: {
  medias: Media[]
  loaded: boolean
  albumNotFound: boolean
}) {
  const {width} = useWindowDimensions()
  const cols = width > xsBreakpoint ? Math.floor((width - drawerWidth) / 200) : 2 // keeps images ~200px wide on non-mobile screen

  const estimatedWidth = width > xsBreakpoint ? (width - drawerWidth) / cols : width / cols
  const estimatedWidthPercentage = estimatedWidth / width

  if (!loaded) {
    return (
      <ImageList cols={cols} gap={2}>
        {Array.from(Array(3 * cols).keys()).map(i => (
          <ImageListItem key={i}>
            <Skeleton variant="rectangular" animation='wave' height={150}/>
          </ImageListItem>
        ))}
      </ImageList>
    )
  }

  if (albumNotFound) {
    return (
      <Alert severity='warning' sx={{mt: 3}}>Album doesn't exist or is not accessible, ask his owner to grant you access
        or choose
        another one on the left.</Alert>
    )
  }
  return (
    <>
      {(medias.length === 0 && (
        <Alert severity='info' sx={{mt: 3}}>This album is empty, uploads some to see them here.</Alert>
      )) || (
        <ImageList cols={cols} gap={0}>
          {/*<ImageListItem key="Subheader" cols={4}>*/}
          {/*  <ListSubheader component="div">December</ListSubheader>*/}
          {/*</ImageListItem>*/}
          {medias.map((media) => (
            <ImageInList media={media} key={media.id} imageMinHeight={estimatedWidth * 9 / 16}
                         imageViewportPercentage={Math.round(estimatedWidthPercentage * 100)}/>
          ))}
        </ImageList>
      )}
    </>
  );
}