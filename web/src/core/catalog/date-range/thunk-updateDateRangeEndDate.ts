import {createSimpleThunkDeclaration} from "src/libs/dthunks";
import {endDateUpdated} from "./action-endDateUpdated";

export const updateDateRangeEndDateDeclaration = createSimpleThunkDeclaration(endDateUpdated);
