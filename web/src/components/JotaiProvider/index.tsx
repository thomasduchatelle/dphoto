// const backendSession = readBackendSession();

// JotaiProvider is passing the Jotai variables from SSR to the client side.
export default function JotaiProvider({children}: { children: React.ReactNode }) {
    // return <Provider>
    //     <JotaiReceiver clientSession={toClientSession(backendSession)}>
    //         {children}
    //     </JotaiReceiver>
    // </Provider>
    // return <JotaiReceiver clientSession={toClientSession(backendSession)}>
    //     {children}
    // </JotaiReceiver>
    return children
}