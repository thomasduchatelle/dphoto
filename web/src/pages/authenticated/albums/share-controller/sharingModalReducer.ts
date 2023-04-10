import {Sharing} from "../../../../core/catalog";

export interface ShareError {
    type: "adding" | "general"
    message: string
}

export interface ShareState {
    open: boolean
    sharedWith: Sharing[]
    error?: ShareError
}

export type OpenSharingModalAction = {
    type: "OpenSharingModalAction"
    sharedWith: Sharing[]
}
export type AddSharingAction = {
    type: "AddSharingAction"
    sharing: Sharing
}
export type RemoveSharingAction = {
    type: "RemoveSharingAction"
    email: string
}
export type CloseSharingModalAction = {
    type: "CloseSharingModalAction"
}
export type SharingModalErrorAction = {
    type: "SharingModalErrorAction"
    error: ShareError
}
export type SharingModalAction =
    OpenSharingModalAction
    | AddSharingAction
    | RemoveSharingAction
    | CloseSharingModalAction
    | SharingModalErrorAction

export function sharingModalReducer(current: ShareState, action: SharingModalAction): ShareState {
    const sortShares = (shares: Sharing[]): Sharing[] => shares.sort((a, b) => a.user.name.toLowerCase() > b.user.name.toLowerCase() ? 1 : -1)

    switch (action.type) {
        case "CloseSharingModalAction":
            return {...current, open: false}

        case "OpenSharingModalAction":
            return {open: true, sharedWith: sortShares(action.sharedWith)}

        case "AddSharingAction":
            const shares = current.sharedWith.filter(s => s.user.email !== action.sharing.user.email);
            shares.push(action.sharing)
            return {open: current.open, sharedWith: sortShares(shares)}

        case "RemoveSharingAction":
            return {open: current.open, sharedWith: current.sharedWith.filter(s => s.user.email !== action.email)}

        case "SharingModalErrorAction":
            return {...current, error: action.error}
    }

    return current
}