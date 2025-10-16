There are some designs issue on how `web/src/components/albums/AlbumsListActions/AlbumListActions.tsx` is taking its arguments and it reduce my ability to make
some updates.

You need to update the code so:

1. `AlbumListActions` arguments are defined as two interfaces: `AlbumListActionsProps` and `AlbumListActionsCallbacks`. The first having the properties derived
   from the state, and the second all the thunks and callbacks.
2. a new selector must returns the same properties as `AlbumListActionsProps` ; remember from the coding standards, the selectors are tested alongside the
   actions that mutate the state. Update the existing tests and add new ones to cover the new selector
3. then you can use the selectors to get the data from the state and the new interface to pass the data along until `AlbumListActions`.

Once you have completed this work, add a new properties to that state called `deleteButtonEnabled` which will be true if at least one of the albums is owned by
the user (false otherwise, even if there is no album). That properties will be used to control if the delete button is enabled or disabled.

For backward compatibility, `AlbumListActions` should have a default value to `false` so the button is per default enabled. And a new story needs to be created
where the button is disabled (visual testing).