'use client';

import {atom, useAtomValue} from "jotai";
import {AccessToken, AuthenticatedUser, ClientSession} from "./security-model";
import {AccessForbiddenError} from "../application";

export const clientSessionAtom = atom<ClientSession | null>(null);

/** Initialized by the JotaiProvider on the client side. */
export const tokenHolder: {
    accessToken?: AccessToken
} = {}

export const useAuthenticatedUser = (): AuthenticatedUser => {
    const clientSession = useAtomValue(clientSessionAtom);
    if (!clientSession) {
        throw new AccessForbiddenError("user is not authenticated")
    }
    return clientSession.authenticatedUser;
}