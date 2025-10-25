'use client';

import {DPhotoApplication} from "./DPhotoApplication";
import {createContext, Dispatch, ReactNode, useReducer} from "react";
import {initialSecurityState, SecurityAction, securityContextReducer, securityContextReducerSupports, SecurityState} from "../security";
import {GeneralState} from "./application-model";
import {ApplicationGenericAction, applicationGenericReducer} from "./application-reducer";
import {Session} from "../../components/AuthProvider";
import {atom, useAtomValue} from "jotai";

export interface ApplicationContextType {
    application: DPhotoApplication

    security: SecurityState

    general: GeneralState
}

const initialAppContext: ApplicationContextTypeWithSetter = {
    context: {
        application: new DPhotoApplication(),
        security: initialSecurityState,
        general: {googleClientId: ''},
    },
    dispatch: () => {
    },
};

export interface ApplicationContextTypeWithSetter {
    context: ApplicationContextType
    dispatch: Dispatch<ApplicationAction>
}


export const securityAtom = atom<Session | undefined>()

export const ApplicationContext = createContext<ApplicationContextTypeWithSetter>(initialAppContext)

export type ApplicationAction = SecurityAction | ApplicationGenericAction

const applicationReducer = (current: ApplicationContextTypeWithSetter, action: ApplicationAction): ApplicationContextTypeWithSetter => {
    let nextContext = current.context

    if (action.type === 'unrecoverable-error') {
        nextContext = applicationGenericReducer(current.context, action)

    } else if (securityContextReducerSupports(action)) {
        nextContext = securityContextReducer(current.context, action);
    }

    if (!Object.is(current.context, nextContext)) {
        return {
            ...current,
            context: nextContext,
        }
    }

    return current
}

const timeoutId = setTimeout(() => {
    console.log("Refreshed hardcoded timeout ticked")

}, 3600)

export const ApplicationContextComponent = ({children, serverSession}: {
    children?: ReactNode
    serverSession: Session
}) => {
    // This is not reading the value magically from the server !
    const securityValue = useAtomValue(securityAtom)
    const [context, dispatch] = useReducer(applicationReducer, applicationReducer({
        ...initialAppContext,
        context: {
            ...initialAppContext.context,
            general: {googleClientId: serverSession.googleClientId},
        }
    }, {
        type: 'authenticated',
        accessToken: serverSession.accessToken,
        user: serverSession.user,
        refreshTimeoutId: timeoutId,
    }))

    return (
        <ApplicationContext.Provider value={{context: context.context, dispatch}}>
            {children}
        </ApplicationContext.Provider>
    );
}
