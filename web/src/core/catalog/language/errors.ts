export interface CatalogError extends Error {
    errorCode: string
    errorMessage?: string
}

export function isCatalogError(err: any): err is CatalogError {
    return err.errorCode
}