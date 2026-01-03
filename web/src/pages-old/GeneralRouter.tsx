'use client';

import AuthenticatedRouter from "./authenticated/AuthenticatedRouter";
import React from "react";
import {useGlobalError} from "../core/application";
import ErrorPage from "./ErrorPage";

const GeneralRouter = () => {
    const globalError = useGlobalError()

    if (globalError) {
        return <ErrorPage error={globalError}/>
    }

    return <AuthenticatedRouter/>
}

export default GeneralRouter
