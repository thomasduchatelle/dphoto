import DownloadIcon from "@mui/icons-material/Download";
import {IconButton, ImageListItem, ImageListItemBar} from "@mui/material";
import {dateTimeToString} from "../../../../core/date-utils";
import {Media} from "../../../../core/domain/catalog";

export function ImageInList({media, imageMinHeight, imageViewportPercentage}: {
  media: Media,
  imageMinHeight: number,
  imageViewportPercentage: number,
}) {
  return (
    <ImageListItem
      key={media.id}
      sx={{
        minHeight: `${imageMinHeight}px`,
        '& .MuiImageListItemBar-root': {
          display: 'none'
        },
        '&:hover .MuiImageListItemBar-root': {
          display: 'flex'
        }
      }}
    >
      <img
        src={`${media.path}`}
        srcSet={`${media.path}?w=150 150w, ${media.path}?w=300 300w, ${media.path}?w=450 450w, ${media.path}?w=800 800w, ${media.path}?w=1200 1200w, ${media.path} 1300w`}
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