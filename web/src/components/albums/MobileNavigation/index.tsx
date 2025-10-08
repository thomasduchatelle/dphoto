'use client';

import HomeIcon from '@mui/icons-material/Home';
import {Breadcrumbs, Link, Typography} from "@mui/material";
import {Album} from "../../../core/catalog";
import {useClientRouter} from "../../ClientRouter";

export default function MobileNavigation({album}: {
    album?: Album
}) {
    const {navigate} = useClientRouter();

    const handleAlbumsClick = (e: React.MouseEvent) => {
        e.preventDefault();
        navigate('/albums');
    };

    return album ? (
        <Breadcrumbs aria-label="breadcrumb">
            <Link underline="hover" color="inherit" href="/albums" onClick={handleAlbumsClick} sx={{display: 'flex', alignItems: 'center', cursor: 'pointer'}}>
                <HomeIcon sx={{mr: 0.5}} fontSize="inherit"/>
                Albums
            </Link>
            <Typography color="text.primary">{album.name}</Typography>
        </Breadcrumbs>
    ) : (
        <Breadcrumbs aria-label="breadcrumb">
            <Typography color="text.primary" sx={{display: 'flex', alignItems: 'center'}}>
                <HomeIcon sx={{mr: 0.5}} fontSize="inherit"/> Albums
            </Typography>
        </Breadcrumbs>
    )
}