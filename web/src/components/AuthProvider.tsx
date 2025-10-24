import {AccessToken, AuthenticatedUser} from "../core/security";
import {atom} from "jotai";
import {getContextData} from "waku/middleware/context";

export interface Session {
    accessToken: AccessToken
    user: AuthenticatedUser
}

const {session} = getContextData() as { session: Session };

console.log(`session loaded with ${JSON.stringify(session)}`)
export const securityAtom = atom<Session>(session)
