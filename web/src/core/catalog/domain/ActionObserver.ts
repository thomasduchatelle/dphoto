export interface HasType {
    type: string
}

export type ActionObserver<T extends HasType> = (action: T) => void