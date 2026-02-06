import {createSimpleThunkDeclaration} from "@/libs/dthunks";
import {endDateUpdated} from "./action-endDateUpdated";

export const updateDateRangeEndDateDeclaration = createSimpleThunkDeclaration(endDateUpdated);
