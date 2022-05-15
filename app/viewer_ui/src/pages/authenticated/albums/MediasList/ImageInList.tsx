import DownloadIcon from "@mui/icons-material/Download";
import {IconButton, ImageListItem, ImageListItemBar} from "@mui/material";
import {dateTimeToString} from "../../../../core/date-utils";
import {Media, MediaType} from "../logic";

export function ImageInList({media, imageViewportPercentage}: {
  media: Media,
  imageViewportPercentage: number,
}) {
  const imageSrc = media.type === MediaType.IMAGE ? `${media.path}` : '/video-placeholder.png';
  const imageSrcSet = media.type === MediaType.IMAGE ? `${media.path}&w=150 150w, ${media.path}&w=300 300w` : '/video-placeholder.png';
  return (
    <ImageListItem
      key={media.id}
      sx={{
        overflow: 'hidden',
        justifyContent: 'center',
        alignItems: 'center',
        '& .MuiImageListItemBar-root': {
          display: 'none'
        },
        '&:hover .MuiImageListItemBar-root': {
          display: 'flex'
        }
      }}
    >
      <img
        src={`${imageSrc}`}
        srcSet={`${imageSrcSet}`}
        sizes={`${imageViewportPercentage}vw`}
        alt={dateTimeToString(media.time)}
        loading="lazy"
      />
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
            href={media.path}
            download
          >
            <DownloadIcon/>
          </IconButton>
        }
      />
    </ImageListItem>
  );
}