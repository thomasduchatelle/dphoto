import LogoutIcon from '@mui/icons-material/Logout';
import PersonIcon from '@mui/icons-material/Person';
import {Avatar, IconButton, ListItemIcon, Menu, MenuItem, Tooltip} from "@mui/material";
import {MouseEvent, useState} from "react";
import {AuthenticatedUser} from "../../core/security";

const UserMenu = ({user, onLogout}: {
    user: AuthenticatedUser,
    onLogout: () => void,
}) => {
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

    const signOut = () => onLogout()

    const handleOpen = (event: MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget)
    }

    const handleClose = () => {
        setAnchorEl(null)
    }

    return (
        <>
            <Tooltip title="Open settings">
                <IconButton sx={{p: 0}} onClick={handleOpen}>
                    <Avatar alt={user.name} src={user.picture ?? "/static/images/avatar/2.jpg"}/>
                </IconButton>
            </Tooltip>
            <Menu
                id="menu-appbar"
                anchorEl={anchorEl}
                anchorOrigin={{
                    vertical: 'bottom',
                    horizontal: 'right',
                }}
                keepMounted
                open={Boolean(anchorEl)}
                onClose={handleClose}
            >
                <MenuItem sx={{cursor: 'unset'}}>
                    <ListItemIcon><PersonIcon/></ListItemIcon>
                    {user.name}
                </MenuItem>
                <MenuItem onClick={signOut}>
                    <ListItemIcon><LogoutIcon/></ListItemIcon>
                    Logout
                </MenuItem>
            </Menu>
        </>
    )
}

export default UserMenu
