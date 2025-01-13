import {CreateAlbumController, CreateAlbumControls, CreateAlbumState, emptyCreateAlbum} from "../../../../core/catalog";
import {FC, useMemo, useState} from "react";
import dayjs, {Dayjs} from "dayjs";
import {CreateAlbumDialog} from "./CreateAlbumDialog";
import {useCatalogAPIAdapter, useCatalogContext} from "../../../../core/catalog-react";

export const CreateAlbumDialogContainer = ({children: Child, firstDay}: { children: FC<CreateAlbumControls>, firstDay?: Dayjs }) => {
    const {handlers} = useCatalogContext()
    const catalogAPIAdapter = useCatalogAPIAdapter()
    const [state, setState] = useState<CreateAlbumState>(emptyCreateAlbum(dayjs()))

    const controller = useMemo(() => new CreateAlbumController(
        setState,
        catalogAPIAdapter,
        handlers,
        firstDay,
    ), [setState, catalogAPIAdapter, handlers, firstDay])


    return <>
        <Child {...controller}/>
        <CreateAlbumDialog state={state} {...controller}/>
    </>
}