import {Alert, ImageList, ImageListItem, ListSubheader, Skeleton} from "@mui/material";
import {Fragment} from "react";
import {toLocaleDateWithDay} from "../../../../core/date-utils";
import useWindowDimensions from "../../../../core/window-utils";
import {MediaWithinADay, mobileBreakpoint} from "../logic";
import {ImageInList} from "./ImageInList";


const drawerWidth = 498;
const marginsWidth = 33;

export default function MediaList({medias, loaded, albumNotFound}: {
  medias: MediaWithinADay[]
  loaded: boolean
  albumNotFound: boolean
}) {
  const {width} = useWindowDimensions()
  // keep images ~200px wide for non-mobile and ~100px for mobile
  const cols = width > mobileBreakpoint ? Math.floor((width - drawerWidth) / 200) : Math.max(Math.floor((width - marginsWidth) / 100), 2)

  const estimatedWidth = width > mobileBreakpoint ? (width - drawerWidth) / cols : (width - marginsWidth) / cols

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
        <Alert severity='info'>This album is empty, uploads some to see them here.</Alert>
      )) || (
        medias.map((day, index) => (
          <Fragment key={day.day.getTime()}>
            <ListSubheader component="div">
              {toLocaleDateWithDay(day.day)}
              {/*{ index === 0 && `- album`}*/}
            </ListSubheader>
            <ImageList cols={cols} gap={0} rowHeight={estimatedWidth}>
              {day.medias.map((media) => (
                <ImageInList media={media} key={media.id} imageViewportPercentage={100 * estimatedWidth / width}/>
              ))}
            </ImageList>
          </Fragment>
        ))
      )}
    </>
  );
}