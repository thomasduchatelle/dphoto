import {Box} from "@mui/material";
import React, {useEffect, useMemo, useState} from "react";
import {useParams} from "react-router-dom";
import {useMustBeAuthenticated} from "../../../core/application";
import {FullHeightLink} from "./FullHeightLink";
import {MediaPageLogic, MediaPageMediasState, MediaPageMediasStateInit} from "./logic";
import MediaNavBar from "./MediaNavBar";

type MediaPageUrlParams = {
  owner: string
  album: string
  encodedId: string
  filename: string
}

export default function MediaPage() {
  const mustBeAuthenticated = useMustBeAuthenticated()
  const {owner, album, encodedId, filename} = useParams<MediaPageUrlParams>()
  const [state, setState] = useState<MediaPageMediasState>(MediaPageMediasStateInit)

  const logic = useMemo(() => new MediaPageLogic(mustBeAuthenticated, setState), [mustBeAuthenticated, setState])

  useEffect(() => {
    logic.findMediasWithCache(owner ?? "", album ?? "").catch(err => console.log(err))
  }, [logic, owner, album])

  const {
    backToAlbumLink,
    imgSrc,
    previousMediaLink,
    previousMediaSrc,
    nextMediaLink,
    nextMediaSrc,
  } = logic.getModel(owner, album, encodedId, filename, state)

  return (
    <>
      {[[previousMediaSrc, '-100%'], [imgSrc, 0], [nextMediaSrc, '100%']].filter(([src, left]) => src && src !== "").map(([src, left]) => (
        <Box
          key={src}
          sx={{
            position: 'absolute',
            top: 0,
            left: left,
            transition: 'left 0.6s',
            backgroundColor: 'black',
            width: '100%',
            height: '100vh',
            display: 'flex',
            alignContent: 'center',
            justifyContent: 'center',
            '& img': {
              maxWidth: '100%',
              maxHeight: '100%',
              objectFit: 'contain',
            }
          }}>
          <img
            src={`${src}&w=1920`}
            srcSet={`${src}&w=360 360w, ${src}&w=1920 1920w`} // TODO 3036 hit the 6MB limits - ${src}&w=3036 3036w, ${src}&w=4048 4048w
            sizes='100vw'
            alt='previous media'/>
        </Box>
      ))}

      <FullHeightLink mediaLink={previousMediaLink} side='left'/>
      <FullHeightLink mediaLink={nextMediaLink} side='right'/>
      <MediaNavBar backUrl={backToAlbumLink} downloadHref={imgSrc}/>
    </>
  )
}