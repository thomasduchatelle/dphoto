import {Box} from "@mui/material";
import React, {useEffect, useMemo, useState} from "react";
import {useNavigate, useParams} from "react-router-dom";
import {useMustBeAuthenticated} from "../../../core/application";
import {FullHeightLink} from "./FullHeightLink";
import {MediaPageLogic, MediaPageMediasState, MediaPageMediasStateInit} from "./logic";
import MediaNavBar from "./MediaNavBar";
import {Key, useNativeControl} from "./useNativeControl";

type MediaPageUrlParams = {
  owner: string
  album: string
  encodedId: string
  filename: string
}

interface ImageRef {
  src: string
  positionLeft: string
  isImage: boolean
}

export default function MediaPage() {
  const mustBeAuthenticated = useMustBeAuthenticated()
  const navigate = useNavigate()
  const {owner, album, encodedId, filename} = useParams<MediaPageUrlParams>()
  const [state, setState] = useState<MediaPageMediasState>(MediaPageMediasStateInit)

  const logic = useMemo(() => new MediaPageLogic(mustBeAuthenticated, setState), [mustBeAuthenticated, setState])

  useEffect(() => {
    if (owner && album) {
      logic.findMediasWithCache(owner, album).catch(err => console.log(err))
    }
  }, [logic, owner, album])

  const {
    backToAlbumLink,
    imgSrc,
    currentIsImage,
    previousMediaLink,
    previousMediaSrc,
    previousIsImage,
    nextMediaLink,
    nextMediaSrc,
    nextIsImage,
  } = logic.getModel(owner, album, encodedId, filename, state)

  const images: ImageRef[] = [
    {
      src: previousMediaSrc,
      positionLeft: '-100%',
      isImage: previousIsImage,
    },
    {
      src: imgSrc,
      positionLeft: '0',
      isImage: currentIsImage,
    },
    {
      src: nextMediaSrc,
      positionLeft: '100%',
      isImage: nextIsImage,
    },
  ]

  const handlers = useNativeControl(key => {
    switch (key) {
      case Key.Left:
        if (previousMediaLink) {
          navigate(previousMediaLink)
        }
        break

      case Key.Right:
        if (nextMediaLink) {
          navigate(nextMediaLink)
        }
        break

      case Key.D:
        document.getElementById("download_current")?.click()
        break

      case Key.Esc:
        document.getElementById("go_back")?.click()
        break

    }
  }, Key.Left, Key.Right, Key.D, Key.Esc)

  return (
    <div {...handlers}>
      {images.filter(ref => ref.src && ref.src !== "")
        .map(({src, positionLeft, isImage}) => (
          <Box
            key={src}
            sx={{
              position: 'absolute',
              top: 0,
              left: positionLeft,
              transition: 'left 0.6s',
              backgroundColor: 'black',
              width: '100%',
              height: '100vh',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              '& img.picture': {
                maxWidth: '100%',
                maxHeight: '100%',
                objectFit: 'contain',
              }
            }}>
            {isImage ? (<img
              className="picture"
              src={`${src}&w=360`}
              srcSet={`${src}&w=360 360w, ${src}&w=1440 1440w, ${src}&w=2400 2400w`}
              sizes='100vw'
              alt='previous media'/>) : (
              <img src="/video-placeholder.png" alt="video placeholder" style={{
                width: '350px',
                height: '288px',
              }}/>
            )}
          </Box>
        ))}

      <FullHeightLink mediaLink={previousMediaLink} side='left'/>
      <FullHeightLink mediaLink={nextMediaLink} side='right'/>
      <MediaNavBar backUrl={backToAlbumLink} backId="go_back" downloadHref={imgSrc} downloadId="download_current"/>
    </div>
  )
}