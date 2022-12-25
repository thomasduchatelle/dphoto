import DownloadIcon from "@mui/icons-material/Download";
import {IconButton, ImageListItem, ImageListItemBar} from "@mui/material";
import {Link} from "react-router-dom";
import {dateTimeToString} from "../../../../core/date-utils";
import {Media, MediaType} from "../logic";

export function ImageInList({media, imageViewportPercentage}: {
  media: Media,
  imageViewportPercentage: number,
}) {
  const imageSrc = media.type === MediaType.IMAGE ? `${media.contentPath}` : '/video-placeholder.png';
  const imageSrcSet = media.type === MediaType.IMAGE ? `${media.contentPath}&w=180 180w, ${media.contentPath}&w=360 360w` : '/video-placeholder.png';
  return (
    <ImageListItem
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