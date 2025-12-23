'use client';

import {DPhotoApplication} from "./DPhotoApplication";
import {createContext, Dispatch, ReactNode, useReducer} from "react";
import {GeneralState} from "./application-model";
import {ApplicationGenericAction, applicationGenericReducer} from "./application-reducer";

export interface ApplicationContextType {
    application: DPhotoApplication

    general: GeneralState
}

const initialAppContext: ApplicationContextTypeWithSetter = {
    context: {
        application: new DPhotoApplication(),
        general: {},
    },
    dispatch: () => {
    },
};

export interface ApplicationContextTypeWithSetter {
    context: ApplicationContextType
    dispatch: Dispatch<ApplicationAction>
}


export const ApplicationContext = createContext<ApplicationContextTypeWithSetter>(initialAppContext)

export type ApplicationAction = ApplicationGenericAction

const applicationReducer = (current: ApplicationContextTypeWithSetter, action: ApplicationAction): ApplicationContextTypeWithSetter => {
    let nextContext = current.context

    if (action.type === 'unrecoverable-error') {
        nextContext = applicationGenericReducer(current.context, action)
    }

    if (!Object.is(current.context, nextContext)) {
        return {
            ...current,
            context: nextContext,
        }
    }

    return current
}

export const ApplicationContextComponent = ({children}: {
    children?: ReactNode
}) => {
    const [context, dispatch] = useReducer(applicationReducer, initialAppContext)

    return (
        <ApplicationContext.Provider value={{context: context.context, dispatch}}>
            {children}
        </ApplicationContext.Provider>
    );
}
