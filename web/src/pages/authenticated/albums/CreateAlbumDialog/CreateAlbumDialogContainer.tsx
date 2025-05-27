import {CreateAlbumController, CreateAlbumControls, CreateAlbumState, CreateAlbumThunk, emptyCreateAlbum} from "../../../../core/catalog";
import {FC, useMemo, useState} from "react";
import dayjs, {Dayjs} from "dayjs";
import {CreateAlbumDialog} from "./CreateAlbumDialog";

export const CreateAlbumDialogContainer = ({createAlbum, children: Child, firstDay}: {
    children: FC<CreateAlbumControls>,
    firstDay?: Dayjs,
    createAlbum: CreateAlbumThunk
}) => {
    const [state, setState] = useState<CreateAlbumState>(emptyCreateAlbum(dayjs()))

    const {openDialogForCreateAlbum, ...controller} = useMemo(() => {
        return new CreateAlbumController(
            setState,
            createAlbum,
            firstDay,
        )
    }, [setState, createAlbum, firstDay])

    return <>
        <Child {...{openDialogForCreateAlbum}}/>
        <CreateAlbumDialog state={state} {...controller}/>
    </>
}