import {CreateAlbumController, CreateAlbumControls, CreateAlbumState, emptyCreateAlbum} from "../../../../core/catalog/domain/CreateAlbumController";
import {FC, useMemo, useState} from "react";
import dayjs, {Dayjs} from "dayjs";
import {CreateAlbumDialog} from "./CreateAlbumDialog";
import {CreateAlbumRequest} from "../../../../core/catalog";

export const CreateAlbumDialogContainer = ({children: Child, firstDay}: { children: FC<CreateAlbumControls>, firstDay?: Dayjs }) => {
    const [state, setState] = useState<CreateAlbumState>(emptyCreateAlbum(dayjs()));

    const controller = useMemo(() => new CreateAlbumController(setState, {
        createAlbum: async (request: CreateAlbumRequest) => {
            console.log("Create album", request);
            return;
        }
    }, firstDay), [setState])


    return <>
        <Child {...controller}/>
        <CreateAlbumDialog state={state} {...controller}/>
    </>
}