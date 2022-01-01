import {GoogleLogout} from "react-google-login";

const LogoutButton = ({clientId, onLogoutSuccess, onLogoutFailure}: {
  clientId: string
  onLogoutSuccess: () => void
  onLogoutFailure?: () => void
}) => (
  <GoogleLogout clientId={clientId} onLogoutSuccess={onLogoutSuccess} onFailure={onLogoutFailure}/>
)

export default LogoutButton