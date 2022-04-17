import AppNavComponent from "../../components/app-nav.component";
import UserMenu from "../../components/user.menu";
import {useSecurityContext, useSignOutCase} from "../../core/application";

const IntegratedNavBar = () => {
  const {user} = useSecurityContext();
  const signOutCase = useSignOutCase()

  if (!user) {
    return null
  }

  return (
    <AppNavComponent
      rightContent={<UserMenu user={user} onLogout={signOutCase.logout}/>}
    />
  );
}

export default IntegratedNavBar
