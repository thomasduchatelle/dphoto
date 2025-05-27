import {Album, AlbumFilterCriterion, AlbumFilterEntry, AlbumsFilteredAction, catalogActions, CatalogViewerAction} from "../domain";
import {onAlbumFilterFunction} from "./thunks-onAlbumFilterChange";

describe('onAlbumFilterChange', () => {

    const albumOtherProps = {totalCount: 0, relativeTemperature: 0, temperature: 0, sharedWith: []}
    const selfOwnedAlbum: Album = {
        albumId: {owner: 'mine', folderName: 'album1'},
        name: 'Jan 25',
        start: new Date(2025, 0, 1),
        end: new Date(2025, 1, 0),
        ...albumOtherProps,
    }
    const someoneElseOwnedAlbum: Album = {
        albumId: {owner: 'someone-else', folderName: 'album2'},
        name: 'Feb 25',
        start: new Date(2025, 1, 1),
        end: new Date(2025, 2, 0),
        ownedBy: {name: 'Someone Else', users: []},
        ...albumOtherProps,
    }

    const selfOwnedFilterEntry: AlbumFilterEntry = {avatars: [], criterion: {selfOwned: true, owners: []}, name: "My albums"}

    const tests: [string, {
        partialState: { selectedAlbum?: Album, allAlbums: Album[] },
        criterion: AlbumFilterCriterion,
        expectedActions: AlbumsFilteredAction[]
    }][] = [
        ["should trigger the AlbumFilteredAction with the new criterion if the selected AlbumID is still selected",
            {
                partialState: {selectedAlbum: selfOwnedAlbum, allAlbums: [selfOwnedAlbum]},
                criterion: selfOwnedFilterEntry.criterion,
                expectedActions: [
                    catalogActions.albumsFilteredAction({criterion: selfOwnedFilterEntry.criterion}),
                ]
            }],
        ["should change the current albumId if the new filter is filtering out the currently selected album",
            {
                partialState: {selectedAlbum: someoneElseOwnedAlbum, allAlbums: [selfOwnedAlbum, someoneElseOwnedAlbum]},
                criterion: selfOwnedFilterEntry.criterion,
                expectedActions: [
                    catalogActions.albumsFilteredAction({criterion: selfOwnedFilterEntry.criterion, redirectTo: selfOwnedAlbum.albumId}),
                ]
            }],
        ["should select an album if no album was selected",
            {
                partialState: {selectedAlbum: undefined, allAlbums: [selfOwnedAlbum, someoneElseOwnedAlbum]},
                criterion: selfOwnedFilterEntry.criterion,
                expectedActions: [
                    catalogActions.albumsFilteredAction({criterion: selfOwnedFilterEntry.criterion, redirectTo: selfOwnedAlbum.albumId}),
                ]
            }],
        ["should not select any album if there is no album in the list",
            {
                partialState: {selectedAlbum: undefined, allAlbums: []},
                criterion: selfOwnedFilterEntry.criterion,
                expectedActions: [
                    catalogActions.albumsFilteredAction({criterion: selfOwnedFilterEntry.criterion}),
                ]
            }],
    ]

    it.each(tests)("%s", (name, {criterion, partialState, expectedActions}) => {
        const collector: CatalogViewerAction[] = []
        const callback = onAlbumFilterFunction.bind(null, collector.push.bind(collector), partialState)

        callback(criterion);

        expect(collector).toEqual(expectedActions)
    })
})

export {}
