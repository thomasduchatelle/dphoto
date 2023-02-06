import AuthenticatedRouter from "./authenticated/AuthenticatedRouter";
import LoginPage from "./Login";
import React, {ReactElement, useCallback, useState} from "react";
import {useIsAuthenticated} from "../core/application/security-hooks";
import {useGlobalError} from "../core/application";
import ErrorPage from "./ErrorPage";
import {Navigate, useSearchParams} from "react-router-dom";

// interface PathParam {
//     path?: string
// }

const RestoreAPIGatewayOriginalPath = ({children}: {
    children: ReactElement;
}) => {
    // note - API Gateway + S3 static will redirect on '/?path=<previously requested url>' when a page is reloaded
    const [search] = useSearchParams();
    const path = search.get('path')
    if (path) {
        return (
            <Navigate to={path}/>
        )
    }

    return children
}

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
        <RestoreAPIGatewayOriginalPath>
            <LoginPage onSuccessfulAuthentication={onSuccessfulAuthentication}/>
        </RestoreAPIGatewayOriginalPath>
    );
}

export default GeneralRouter
