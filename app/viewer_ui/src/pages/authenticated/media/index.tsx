import {Box} from "@mui/material";
import React from "react";
import {useParams} from "react-router-dom";
import MediaNavBar from "./MediaNavBar";

type MediaPageUrlParams = {
  owner: string
  album: string
  encodedId: string
  filename: string
}

export default function MediaPage() {
  const {owner, album, encodedId, filename} = useParams<MediaPageUrlParams>()
  const imgSrc = `/api/v1/owners/${owner}/medias/${encodedId}/${filename}`

  return (
    <>
      <MediaNavBar backUrl={`/albums/${owner}/${album}`} downloadHref={imgSrc}/>
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
          src={`${imgSrc}?w=1920`}
          srcSet={`${imgSrc}?w=360 360w, ${imgSrc}?w=1920 1920w, ${imgSrc}?w=3036 3036w, ${imgSrc}?w=4048 4048w`}
          sizes='100vw'
          alt={filename}/>
      </Box>
    </>
  )
}