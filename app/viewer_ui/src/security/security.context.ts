import {createContext} from "react";
import {SecurityContextType} from "./security.model";

export const SecurityReactContext = createContext<SecurityContextType>({loggedUser: undefined})
