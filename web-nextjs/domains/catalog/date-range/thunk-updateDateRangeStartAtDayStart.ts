import {createSimpleThunkDeclaration} from "@/libs/dthunks";
import {atStartDayToggled} from "./action-atStartDayToggled";

export const updateDateRangeStartAtDayStartDeclaration = createSimpleThunkDeclaration(atStartDayToggled);
