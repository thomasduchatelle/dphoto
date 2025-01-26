import {CreateAlbumController, CreateAlbumControls, CreateAlbumState, emptyCreateAlbum} from "../../../../core/catalog/domain/CreateAlbumController";
import {FC, useMemo, useState} from "react";
import dayjs from "dayjs";
import CreateAlbumDialog from "./CreateAlbumDialog";

// export type CreateAlbumDialogContainerChild = (props: { state: CreateAlbumState } & CreateAlbumControls) => JSX.Element;
export type CreateAlbumDialogChildProps = { state: CreateAlbumState } & CreateAlbumControls;
export type CreateAlbumDialogContainerChild = FC<CreateAlbumDialogChildProps>;

export const CreateAlbumDialogContainer = ({foo: Child}: { foo: CreateAlbumDialogContainerChild }) => {
    const [state, setState] = useState<CreateAlbumState>(emptyCreateAlbum(dayjs()));

    const controller = useMemo(() => new CreateAlbumController(setState, {
        createAlbum: async () => {
            return;
        }
    }), [setState])


    return <>
        <Child state={state} {...controller}/>
        <CreateAlbumDialog state={state} {...controller}/>
    </>
}

export const Foo = () => {
    return <CreateAlbumDialogContainer foo={
        (props: CreateAlbumDialogChildProps) => {
            return <div>
                <h1>Foo</h1>
                <button onClick={props.openNew}>Open</button>
            </div>
        }}/>
    // return <CreateAlbumDialogContainer>
    //     {(props: CreateAlbumDialogChildProps) => {
    //         return <div>
    //             <h1>Foo</h1>
    //             <button onClick={props.openNew}>Open</button>
    //         </div>
    //     }}
    //     {/*{({props: {state}}): JSX.Element => {*/}
    //     {/*    return <div>*/}
    //     {/*        <h1>Foo</h1>*/}
    //     {/*        <button onClick={props.openNew}>Open</button>*/}
    //     {/*    </div>*/}
    //     {/*}}*/}
    // </CreateAlbumDialogContainer>
}