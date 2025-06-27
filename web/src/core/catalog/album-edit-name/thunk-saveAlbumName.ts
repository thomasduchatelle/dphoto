import {ThunkDeclaration} from "src/libs/dthunks";
import {AlbumId, CatalogViewerState, getErrorMessage, isCatalogError} from "../language";
import {albumRenamingStarted, AlbumRenamingStarted} from "./action-albumRenamingStarted";
import {albumRenamed, AlbumRenamed} from "./action-albumRenamed";
import {albumRenamingFailed, AlbumRenamingFailed} from "./action-albumRenamingFailed";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {CatalogAPIAdapter} from "../adapters/api";

export interface SaveAlbumNamePort {
    renameAlbum(albumId: AlbumId, newName: string, newFolderName?: string): Promise<AlbumId>;
}

export interface SaveAlbumNamePreselection {
    albumId: AlbumId;
    albumName: string;
    customFolderName: string;
    isCustomFolderNameEnabled: boolean;
}

export async function saveAlbumNameThunk(
    dispatch: (action: AlbumRenamingStarted | AlbumRenamed | AlbumRenamingFailed) => void,
    saveAlbumNamePort: SaveAlbumNamePort,
    preselection: SaveAlbumNamePreselection
): Promise<void> {
    dispatch(albumRenamingStarted());

    try {
        const newFolderName = preselection.isCustomFolderNameEnabled ? preselection.customFolderName : undefined;
        const newAlbumId = await saveAlbumNamePort.renameAlbum(preselection.albumId, preselection.albumName, newFolderName);
        dispatch(albumRenamed({previousAlbumId: preselection.albumId, newAlbumId, newName: preselection.albumName}));
    } catch (err) {
        if (isCatalogError(err)) {
            dispatch(albumRenamingFailed(err));
        } else {
            dispatch(albumRenamingFailed({message: getErrorMessage(err) || "Something went wrong. Please try again."}));
        }
    }
}

export const saveAlbumNameDeclaration: ThunkDeclaration<
    CatalogViewerState,
    SaveAlbumNamePreselection,
    () => Promise<void>,
    CatalogFactoryArgs
> = {
    selector: (state: CatalogViewerState): SaveAlbumNamePreselection => {
        if (!state.dialog || state.dialog.type !== "EditNameDialog") {
            return {
                albumId: {owner: "", folderName: ""},
                albumName: "",
                customFolderName: "",
                isCustomFolderNameEnabled: false,
            };
        }
        return {
            albumId: state.dialog.albumId,
            albumName: state.dialog.albumName,
            customFolderName: state.dialog.customFolderName,
            isCustomFolderNameEnabled: state.dialog.isCustomFolderNameEnabled,
        };
    },

    factory: ({dispatch, app, partialState}) => {
        const catalogAdapter = new CatalogAPIAdapter(app.axiosInstance, app);
        return saveAlbumNameThunk.bind(null, dispatch, catalogAdapter, partialState);
    },
};
