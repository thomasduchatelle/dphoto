import {AccessToken, AuthenticatedUser} from "../core/security";
import {getContextData} from "waku/middleware/context";
import {ReactNode} from "react";
import {ApplicationContextComponent} from "../core/application";

export interface Session {
    accessToken: AccessToken
    user: AuthenticatedUser
}

const {session} = getContextData() as { session: Session };

export const AuthProvider = ({children}: { children: ReactNode }) => {
    return <ApplicationContextComponent serverSession={session}>{children}</ApplicationContextComponent>
}