import {createSimpleThunkDeclaration} from "@/libs/dthunks";
import {startDateUpdated} from "./action-startDateUpdated";

export const updateDateRangeStartDateDeclaration = createSimpleThunkDeclaration(startDateUpdated);
