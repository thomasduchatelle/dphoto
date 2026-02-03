import {createAction} from "@/libs/daction";
import {CatalogViewerState} from "../language";

export const createDialogOpened = createAction<CatalogViewerState>(
    "CreateDialogOpened",
    (current: CatalogViewerState) => {
        const today = new Date();
        const dayOfWeek = today.getDay(); // 0 for Sunday, 6 for Saturday

        // Calculate the Saturday of the previous week
        const startDate = new Date(today);
        startDate.setDate(today.getDate() - (dayOfWeek + 1) % 7 - 7); // Go back to last Saturday
        startDate.setHours(0, 0, 0, 0);

        // Calculate the Monday of the current week
        const endDate = new Date(today);
        endDate.setDate(today.getDate() - (dayOfWeek + 6) % 7); // Go back to this Monday
        endDate.setHours(0, 0, 0, 0);

        return {
            ...current,
            dialog: {
                type: "CreateDialog",
                albumId: {owner: "", folderName: ""}, // Will be set when creating
                albumName: "",
                customFolderName: "",
                isCustomFolderNameEnabled: false,
                nameError: {},
                startDate: startDate,
                endDate: endDate,
                startAtDayStart: true,
                endAtDayEnd: true,
                isLoading: false,
            },
        };
    }
);

export type CreateDialogOpened = ReturnType<typeof createDialogOpened>;
