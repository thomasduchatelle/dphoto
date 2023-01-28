import AuthenticatedRouter from "./authenticated/AuthenticatedRouter";
import LoginPage from "./Login";
import {useCallback, useState} from "react";
import {useIsAuthenticated} from "../core/application/security-hooks";
import {useGlobalError} from "../core/application";
import ErrorPage from "./ErrorPage";

const GeneralRouter = () => {
    const isAuthenticated = useIsAuthenticated();
    const [displayLoginPage, setDisplayLoginPage] = useState<boolean>(!isAuthenticated) // keep loading page displayed while loading on the background
    const onSuccessfulAuthentication = useCallback(() => setDisplayLoginPage(false), [setDisplayLoginPage])
    const globalError = useGlobalError()

    if (globalError) {
        return <ErrorPage error={globalError}/>
    }

    return isAuthenticated && !displayLoginPage ? (
        <AuthenticatedRouter/>
    ) : (
        <LoginPage onSuccessfulAuthentication={onSuccessfulAuthentication}/>
    );
}

export default GeneralRouter
