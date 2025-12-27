'use client';

import DownloadIcon from "@mui/icons-material/Download";
import {IconButton, ImageListItem, ImageListItemBar} from "@mui/material";
import {dateTimeToString} from "../../../core/utils/date-utils";
import {Media, MediaType} from "../../../core/catalog";
import {useEffect, useRef} from "react";
import {useClientRouter} from "../../ClientRouter";

function addQueryParam(url: string, param: string, value: string): string {
    const trimmedUrl = url.replace(/[?&]$/, '');
    const separator = trimmedUrl.includes('?') ? '&' : '?';
    return `${trimmedUrl}${separator}${param}=${value}`;
}

export function ImageInList({media, imageViewportPercentage, autoFocus = false}: {
    media: Media,
    imageViewportPercentage: number,
    autoFocus?: boolean,
}) {
    const {navigate} = useClientRouter();
    const itemRef = useRef<HTMLLIElement | null>(null)
    const imageSrc = media.type === MediaType.IMAGE ? `${media.contentPath}` : '/video-placeholder.png';
    const imageSrcSet = media.type === MediaType.IMAGE ? `${addQueryParam(media.contentPath, 'w', '180')} 180w, ${addQueryParam(media.contentPath, 'w', '360')} 360w` : '/video-placeholder.png';

    useEffect(() => {
        if (autoFocus && itemRef.current) {
            itemRef.current.scrollIntoView({behavior: 'smooth', block: 'start'})
        }
    }, [autoFocus, itemRef])

    const handleClick = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();
        navigate(media.uiRelativePath);
    };

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
            <a href={media.uiRelativePath} onClick={handleClick}>
                <img
                    src={`${imageSrc}`}
                    srcSet={`${imageSrcSet}`}
                    sizes={`${imageViewportPercentage}vw`}
                    alt={dateTimeToString(media.time)}
                    loading="lazy"
                />
            </a>
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