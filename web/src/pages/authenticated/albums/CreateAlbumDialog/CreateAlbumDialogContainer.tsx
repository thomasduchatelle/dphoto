import {AlbumId, CreateAlbumController, CreateAlbumControls, CreateAlbumRequest, CreateAlbumState, emptyCreateAlbum} from "../../../../core/catalog";
import {FC, useMemo, useState} from "react";
import dayjs, {Dayjs} from "dayjs";
import {CreateAlbumDialog} from "./CreateAlbumDialog";

export const CreateAlbumDialogContainer = ({children: Child, firstDay}: { children: FC<CreateAlbumControls>, firstDay?: Dayjs }) => {
    const [state, setState] = useState<CreateAlbumState>(emptyCreateAlbum(dayjs()));

    const controller = useMemo(() => new CreateAlbumController(
        setState,
        {
            createAlbum: async (request: CreateAlbumRequest): Promise<AlbumId> => {
                console.log("Create album", request);
                return {owner: "owner1", folderName: "/album1"}
            }
        },
        {
            onAlbumCreated(albumId: AlbumId): Promise<void> {
                return Promise.resolve();
            }
        },
        firstDay), [setState])


    return <>
        <Child {...controller}/>
        <CreateAlbumDialog state={state} {...controller}/>
    </>
}