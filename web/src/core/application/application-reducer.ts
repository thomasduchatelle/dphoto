import {ApplicationContextType} from "./application-context";

export type ConfigLoadedAction = {
    type: 'config-loaded'
    googleClientId: string
}

export type UnrecoverableErrorAction = {
    type: 'unrecoverable-error'
    error: Error
}

export type ApplicationGenericAction = ConfigLoadedAction | UnrecoverableErrorAction

export function applicationGenericReducer(current: ApplicationContextType, action: ApplicationGenericAction): ApplicationContextType {
    switch (action.type) {
        case 'config-loaded':
            return {
                ...current,
                general: {
                    ...current.general,
                    googleClientId: action.googleClientId,
                }
            }

        case "unrecoverable-error":
            return {
                ...current,
                general: {
                    ...current.general,
                    error: action.error,
                }
            }
    }

    return current
}