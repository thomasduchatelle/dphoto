import {CreateAlbumController, CreateAlbumControls, CreateAlbumState, emptyCreateAlbum} from "../../../../core/catalog/domain/CreateAlbumController";
import {FC, useMemo, useState} from "react";
import dayjs, {Dayjs} from "dayjs";
import {CreateAlbumDialog} from "./CreateAlbumDialog";

export const CreateAlbumDialogContainer = ({children: Child, firstDay}: { children: FC<CreateAlbumControls>, firstDay?: Dayjs }) => {
    const [state, setState] = useState<CreateAlbumState>(emptyCreateAlbum(dayjs()));

    const controller = useMemo(() => new CreateAlbumController(setState, {
        createAlbum: async () => {
            return;
        }
    }, firstDay), [setState])


    return <>
        <Child {...controller}/>
        <CreateAlbumDialog state={state} {...controller}/>
    </>
}