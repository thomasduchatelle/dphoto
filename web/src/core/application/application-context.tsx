import {DPhotoApplication} from "./DPhotoApplication";
import {createContext, Dispatch, ReactNode, useEffect, useReducer} from "react";
import {
    initialSecurityState,
    SecurityAction,
    securityContextReducer,
    securityContextReducerSupports,
    SecurityState
} from "../security";
import {GeneralState} from "./application-model";
import axios from "axios";
import {ApplicationGenericAction, applicationGenericReducer} from "./application-reducer";

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

export const ApplicationContext = createContext<ApplicationContextTypeWithSetter>(initialAppContext)

export type ApplicationAction = SecurityAction | ApplicationGenericAction

const applicationReducer = (current: ApplicationContextTypeWithSetter, action: ApplicationAction): ApplicationContextTypeWithSetter => {
    let nextContext = current.context

    if (action.type === 'config-loaded' || action.type === 'unrecoverable-error') {
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

interface ConfigFile {
    googleClientId: string
}

export const ApplicationContextComponent = ({children}: {
    children?: ReactNode
}) => {
    const [context, dispatch] = useReducer(applicationReducer, initialAppContext)

    useEffect(() => {
        axios.get<ConfigFile>("/env-config.json")
            .then(cfg => {
                dispatch({type: 'config-loaded', googleClientId: cfg.data.googleClientId})
            })
    }, [])

    return (
        <ApplicationContext.Provider value={{context: context.context, dispatch}}>
            {children}
        </ApplicationContext.Provider>
    );
}
