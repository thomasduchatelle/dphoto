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
  }, [owner, album])

  const {
    backToAlbumLink,
    imgSrc,
    previousMediaLink,
    nextMediaLink
  } = logic.getModel(owner, album, encodedId, filename, state)

  console.log(`loading image ${imgSrc}`)

  return (
    <>
      <MediaNavBar backUrl={backToAlbumLink} downloadHref={imgSrc}/>
      <FullHeightLink mediaLink={previousMediaLink} side='left'/>
      <FullHeightLink mediaLink={nextMediaLink} side='right'/>

      <Box sx={{
        backgroundColor: 'black',
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
          key={imgSrc}
          src={`${imgSrc}?w=1920`}
          srcSet={`${imgSrc}?w=360 360w, ${imgSrc}?w=1920 1920w, ${imgSrc}?w=3036 3036w, ${imgSrc}?w=4048 4048w`}
          sizes='100vw'
          alt={filename}/>
      </Box>
    </>
  )
}