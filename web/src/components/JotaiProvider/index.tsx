import {Provider} from "jotai";
import {JotaiReceiver} from "./JotaiReceiver";
import {readBackendSession} from "../../core/security-ssr";
import {toClientSession} from "../../core/security";

const backendSession = readBackendSession();

/** JotaiProvider is passing the Jotai variables from SSR to the client side. */
export default function JotaiProvider({children}: { children: React.ReactNode }) {
    return <Provider>
        <JotaiReceiver clientSession={toClientSession(backendSession)}>
            {children}
        </JotaiReceiver>
    </Provider>
}