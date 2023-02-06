export abstract class ErrorWithPublicMessage extends Error {
    public readonly publicMessage: string = ""
}

export interface GeneralState {
    googleClientId: string
    error?: ErrorWithPublicMessage
}

export interface AccessTokenHolder {
    getAccessToken(): string
}