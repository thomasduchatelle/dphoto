import {createSimpleThunkDeclaration} from "src/libs/dthunks";
import {atEndDayToggled} from "./action-atEndDayToggled";

export const updateDateRangeEndAtDayEndDeclaration = createSimpleThunkDeclaration(atEndDayToggled);
