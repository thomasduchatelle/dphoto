import {createSimpleThunkDeclaration} from "@/libs/dthunks";
import {atEndDayToggled} from "./action-atEndDayToggled";

export const updateDateRangeEndAtDayEndDeclaration = createSimpleThunkDeclaration(atEndDayToggled);
