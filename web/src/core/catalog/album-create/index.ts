import {openCreateDialogDeclaration} from "./thunk-openCreateDialog";
import {closeCreateDialogDeclaration} from "./thunk-closeCreateDialog";
import {updateCreateDialogNameDeclaration} from "./thunk-updateCreateDialogName";
import {updateCreateDialogFolderNameDeclaration} from "./thunk-updateCreateDialogFolderName";
import {updateCreateDialogWithCustomFolderNameDeclaration} from "./thunk-updateCreateDialogWithCustomFolderName";
import {submitCreateAlbumDeclaration} from "./thunk-submitCreateAlbum";
import {createAlbumDeclaration} from "./album-createAlbum";
import {dateRangeThunks} from "../date-range";

export * from "./album-createAlbum";
export * from "./selector-createDialogSelector";

/**
 * Thunks related to album creation.
 *
 * Expected handler types:
 * - `openCreateDialog`: `() => void`
 * - `closeCreateDialog`: `() => void`
 * - `updateCreateDialogName`: `(name: string) => void`
 * - `updateCreateDialogStartDate`: `(date: Date | null) => void`
 * - `updateCreateDialogEndDate`: `(date: Date | null) => void`
 * - `updateCreateDialogFolderName`: `(folderName: string) => void`
 * - `updateCreateDialogWithCustomFolderName`: `(withCustom: boolean) => void`
 * - `updateCreateDialogStartsAtStartOfTheDay`: `(startsAtStart: boolean) => void`
 * - `updateCreateDialogEndsAtEndOfTheDay`: `(endsAtEnd: boolean) => void`
 * - `submitCreateAlbum`: `() => Promise<void>`
 * - `createAlbum`: `(request: CreateAlbumRequest) => Promise<AlbumId>`
 */
export const albumCreateThunks = {
    openCreateDialog: openCreateDialogDeclaration,
    closeCreateDialog: closeCreateDialogDeclaration,
    updateCreateDialogName: updateCreateDialogNameDeclaration,
    updateCreateDialogStartDate: dateRangeThunks.updateDateRangeStartDate,
    updateCreateDialogEndDate: dateRangeThunks.updateDateRangeEndDate,
    updateCreateDialogFolderName: updateCreateDialogFolderNameDeclaration,
    updateCreateDialogWithCustomFolderName: updateCreateDialogWithCustomFolderNameDeclaration,
    updateCreateDialogStartsAtStartOfTheDay: dateRangeThunks.updateDateRangeStartAtDayStart,
    updateCreateDialogEndsAtEndOfTheDay: dateRangeThunks.updateDateRangeEndAtDayEnd,
    submitCreateAlbum: submitCreateAlbumDeclaration,
    createAlbum: createAlbumDeclaration,
};
