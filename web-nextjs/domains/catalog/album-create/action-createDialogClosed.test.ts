import {CatalogViewerState} from "../language";
import {createDialogClosed} from "./action-createDialogClosed";
import {createDialogPrefilledForMar25, loadedStateWithTwoAlbums} from "../tests/test-helper-state";

describe('action:createDialogClosed', () => {
    it('should close the dialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: createDialogPrefilledForMar25,
        };

        const action = createDialogClosed();
        const newState = action.reducer(state, action);

        expect(newState.dialog).toBeUndefined();
    });
});
