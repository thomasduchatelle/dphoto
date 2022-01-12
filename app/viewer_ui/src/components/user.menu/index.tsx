import {Avatar, IconButton, Menu, MenuItem, Tooltip} from "@mui/material";
import {MouseEvent, useState} from "react";
import {useGoogleLogout} from "react-google-login";
import googleConfig from "../../config/google.config";
import {AuthenticatedUser} from "../../core/domain/security";

export default ({user, onLogout}: {
  user: AuthenticatedUser,
  onLogout: () => void,
}) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const {signOut} = useGoogleLogout({
    clientId: googleConfig.clientId,
    uxMode: googleConfig.uxMode,
    onLogoutSuccess: () => {
      onLogout()
    }
  })

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
        <MenuItem onClick={signOut}>Logout</MenuItem>
      </Menu>
    </>
  )
}