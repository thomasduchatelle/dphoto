'use client';

import {ClientSession} from "../../core/security/security-model";
import {clientSessionAtom, tokenHolder} from "../../core/security/client-utils";
import {useHydrateAtoms} from 'jotai/utils'

export interface JotaiReceiver {
    children?: React.ReactNode
    clientSession: ClientSession | null
}

export const JotaiReceiver = ({children, clientSession}: JotaiReceiver) => {
    useHydrateAtoms([
        [clientSessionAtom, clientSession]
    ])
    if (clientSession?.accessToken) {
        tokenHolder.accessToken = clientSession.accessToken;
    }

    return children;
}