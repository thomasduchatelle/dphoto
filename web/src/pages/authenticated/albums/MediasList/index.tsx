import {Alert, Box, CircularProgress, ImageList, Typography, useMediaQuery, useTheme} from "@mui/material";
import {Fragment} from "react";
import {toLocaleDateWithDay} from "../../../../core/utils/date-utils";
import useWindowDimensions from "../../../../core/utils/window-utils";
import {ImageInList} from "./ImageInList";
import {MediaWithinADay} from "../../../../core/catalog";

const drawerWidth = 498;
const marginsWidth = 33;

export default function MediaList({medias, loaded, albumNotFound, scrollToMedia = undefined}: {
    medias: MediaWithinADay[]
    loaded: boolean
    albumNotFound: boolean
    scrollToMedia?: string
}) {
    const {width} = useWindowDimensions()

    const theme = useTheme()
    const isMobileDevice = useMediaQuery(theme.breakpoints.down('md'));

    // keep images ~200px wide for non-mobile and ~100px for mobile
    const cols = isMobileDevice ? Math.max(Math.floor((width - marginsWidth) / 100), 2) : Math.floor((width - drawerWidth) / 200)

    const estimatedWidth = isMobileDevice ? (width - marginsWidth) / cols : (width - drawerWidth) / cols

    if (!loaded) {
        return (
            <Box sx={{
                paddingTop: 10,
                textAlign: 'center',
            }}>
                <CircularProgress size={150}/>
            </Box>
        )
    }

    if (albumNotFound) {
        return (
            <Alert severity='warning' sx={{mt: 3}}>
                Album doesn't exist or is not accessible, ask his owner to grant you access or choose another one on the
                left.
            </Alert>
        )
    }
    return (
        <>
            {(medias.length === 0 && (
                <Alert severity='info'>This album is empty, uploads some to see them here.</Alert>
            )) || (
                medias.map((day, index) => (
                    <Fragment key={day.day.getTime()}>
                        <Typography variant="h6" sx={theme => ({
                            // fontSize: '28px',
                            fontWeight: 'normal',
                            margin: theme.spacing(3, 0, 1, 0),
                            '&:first-of-type': {
                                marginTop: theme.spacing(1),
                            }
                        })}>
                            {toLocaleDateWithDay(day.day)}
                        </Typography>
                        <ImageList cols={cols} gap={0} rowHeight={estimatedWidth}>
                            {day.medias.map((media) => (
                                <ImageInList
                                    media={media}
                                    key={media.id}
                                    imageViewportPercentage={100 * estimatedWidth / width}
                                    autoFocus={media.id === scrollToMedia}
                                />
                            ))}
                        </ImageList>
                    </Fragment>
                ))
            )}
        </>
    );
}