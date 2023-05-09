import DownloadIcon from "@mui/icons-material/Download";
import {IconButton, ImageListItem, ImageListItemBar} from "@mui/material";
import {Link} from "react-router-dom";
import {dateTimeToString} from "../../../../core/utils/date-utils";
import {Media, MediaType} from "../../../../core/catalog";
import {useEffect, useRef} from "react";

export function ImageInList({media, imageViewportPercentage, autoFocus = false}: {
    media: Media,
    imageViewportPercentage: number,
    autoFocus?: boolean,
}) {
    const itemRef = useRef<HTMLLIElement | null>(null)
    const imageSrc = media.type === MediaType.IMAGE ? `${media.contentPath}` : '/video-placeholder.png';
    const imageSrcSet = media.type === MediaType.IMAGE ? `${media.contentPath}&w=180 180w, ${media.contentPath}&w=360 360w` : '/video-placeholder.png';

    useEffect(() => {
        if (autoFocus && itemRef.current) {
            itemRef.current.scrollIntoView({ behavior: 'smooth', block: 'start' })
        }
    }, [autoFocus, itemRef])

    return (
        <ImageListItem
            ref={itemRef}
            sx={{
                overflow: 'hidden',
                '& .MuiImageListItemBar-root': {
                    display: 'none'
                },
                '&:hover .MuiImageListItemBar-root': {
                    display: {lg: 'flex'}
                },
                '& a': {
                    height: '100%',
                    width: '100%',
                    display: 'flex',
                    justifyContent: 'center',
                    alignItems: 'center',
                    overflow: 'hidden',
                },
                '& img': {
                    minHeight: '100%',
                    minWidth: '100%',
                    objectFit: 'cover',
                },
            }}
        >
            <Link to={`${media.uiRelativePath}`}>
                <img
                    src={`${imageSrc}`}
                    srcSet={`${imageSrcSet}`}
                    sizes={`${imageViewportPercentage}vw`}
                    alt={dateTimeToString(media.time)}
                    loading="lazy"
                />
            </Link>
            <ImageListItemBar
                title={dateTimeToString(media.time)}
                subtitle={media.source ? `@${media.source}` : ''}
                sx={{
                    '& .MuiImageListItemBar-title': {
                        fontSize: '0.8em'
                    },
                    '& .MuiImageListItemBar-subtitle': {
                        fontSize: '0.6em'
                    }
                }}
                actionIcon={
                    <IconButton
                        sx={{color: 'rgba(255, 255, 255, 0.54)'}}
                        aria-label='download image'
                        title='Download'
                        href={media.contentPath}
                        component='a'
                        download
                    >
                        <DownloadIcon/>
                    </IconButton>
                }
            />
        </ImageListItem>
    );
}