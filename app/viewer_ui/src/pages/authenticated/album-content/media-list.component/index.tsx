import {Alert, ImageList, ImageListItem, Skeleton} from "@mui/material";
import {Media} from "../../../../core/domain/catalog";
import useWindowDimensions from "../../../../core/window-utils";
import {ImageInList} from "./media-in-list";


const drawerWidth = 498;
const marginsWidth = 33;
const breakpoint = 1200;

export function MediaListComponent({medias, loaded, albumNotFound}: {
  medias: Media[]
  loaded: boolean
  albumNotFound: boolean
}) {
  const {width} = useWindowDimensions()
  // keep images ~200px wide for non-mobile and ~100px for mobile
  const cols = width > breakpoint ? Math.floor((width - drawerWidth) / 200) : Math.max(Math.floor((width - marginsWidth) / 100), 2)

  const estimatedWidth = width > breakpoint ? (width - drawerWidth) / cols : (width - marginsWidth) / cols
  console.log(`estimatedWidth = (${width} - ${drawerWidth}) / ${cols} = ${estimatedWidth}`)

  if (!loaded) {
    return (
      <ImageList cols={cols} gap={2} rowHeight={estimatedWidth}>
        {Array.from(Array(3 * cols).keys()).map(i => (
          <ImageListItem key={i}>
            <Skeleton variant="rectangular" animation='wave' height={estimatedWidth}/>
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
        <ImageList cols={cols} gap={0} rowHeight={estimatedWidth}>
          {/*<ImageListItem key="Subheader" cols={4}>*/}
          {/*  <ListSubheader component="div">December</ListSubheader>*/}
          {/*</ImageListItem>*/}
          {medias.map((media) => (
            <ImageInList media={media} key={media.id} imageViewportPercentage={100 * estimatedWidth / width}/>
          ))}
        </ImageList>
      )}
    </>
  );
}