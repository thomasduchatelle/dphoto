import HomeIcon from '@mui/icons-material/Home';
import {Breadcrumbs, Link, Typography} from "@mui/material";
import {Album} from "../logic";

export default function MobileNavigation({album}: {
  album?: Album
}) {
  return (album && (
    <Breadcrumbs aria-label="breadcrumb">
      <Link underline="hover" color="inherit" href="/albums" sx={{display: 'flex', alignItems: 'center'}}>
        <HomeIcon sx={{mr: 0.5}} fontSize="inherit"/>
        Albums
      </Link>
      <Typography color="text.primary">{album.name}</Typography>
    </Breadcrumbs>
  )) || (
    <Breadcrumbs aria-label="breadcrumb">
      <Typography color="text.primary" sx={{display: 'flex', alignItems: 'center'}}>
        <HomeIcon sx={{mr: 0.5}} fontSize="inherit"/> Albums
      </Typography>
    </Breadcrumbs>
  )
}