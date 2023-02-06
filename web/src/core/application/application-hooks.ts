import {Dispatch, useContext} from "react";
import {ApplicationContext} from "./application-context";
import {AxiosInstance} from "axios";
import {GeneralState} from "./application-model";
import {UnrecoverableErrorAction} from "./application-reducer";

export const useApplication = () => {
    return useContext(ApplicationContext).context.application
}

export const useAxios = (): AxiosInstance => {
    return useApplication().axiosInstance
}

export const useConfigContext = (): GeneralState => {
    return useContext(ApplicationContext).context.general
}

export const useGlobalError = (): Error | undefined => {
    return useContext(ApplicationContext).context.general.error
}

export const useUnrecoverableErrorDispatch = (): Dispatch<UnrecoverableErrorAction> => {
    return useContext(ApplicationContext).dispatch
}
