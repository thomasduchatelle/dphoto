import {openCreateDialogDeclaration} from "./thunk-openCreateDialog";
import {closeCreateDialogDeclaration} from "./thunk-closeCreateDialog";
import {submitCreateAlbumDeclaration} from "./thunk-submitCreateAlbum";
import {dateRangeThunks} from "../date-range";
import {baseNameEditThunks} from "../base-name-edit";

export * from "./selector-createDialogSelector";

/**
 * Thunks related to album creation.
 *
 * Expected handler types:
 * - `openCreateDialog`: `() => void`
 * - `closeCreateDialog`: `() => void`
 * - `changeAlbumName`: `(albumName: string) => void`
 * - `updateCreateDialogStartDate`: `(date: Date | null) => void`
 * - `updateCreateDialogEndDate`: `(date: Date | null) => void`
 * - `changeFolderName`: `(folderName: string) => void`
 * - `changeFolderNameEnabled`: `(isFolderNameEnabled: boolean) => void`
 * - `updateCreateDialogStartsAtStartOfTheDay`: `(startsAtStart: boolean) => void`
 * - `updateCreateDialogEndsAtEndOfTheDay`: `(endsAtEnd: boolean) => void`
 * - `submitCreateAlbum`: `() => Promise<void>`
 */
export const albumCreateThunks = {
    openCreateDialog: openCreateDialogDeclaration,
    closeCreateDialog: closeCreateDialogDeclaration,
    changeAlbumName: baseNameEditThunks.changeAlbumName,
    updateCreateDialogStartDate: dateRangeThunks.updateDateRangeStartDate,
    updateCreateDialogEndDate: dateRangeThunks.updateDateRangeEndDate,
    changeFolderName: baseNameEditThunks.changeFolderName,
    changeFolderNameEnabled: baseNameEditThunks.changeFolderNameEnabled,
    updateCreateDialogStartsAtStartOfTheDay: dateRangeThunks.updateDateRangeStartAtDayStart,
    updateCreateDialogEndsAtEndOfTheDay: dateRangeThunks.updateDateRangeEndAtDayEnd,
    submitCreateAlbum: submitCreateAlbumDeclaration,
};