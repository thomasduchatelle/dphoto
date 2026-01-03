export abstract class ErrorWithPublicMessage extends Error {
    public readonly publicMessage: string = ""
}

export interface GeneralState {
    error?: ErrorWithPublicMessage
}
