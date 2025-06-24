* [X] Support OrphanedMediasError natively in the UI (edit dates)
* [X] Reuse the navigation actions to reload the medias and albums (but it needs to close the dialog)
* [X] Move the grouping by days of the medias in the action, not in the port
* [X] Any action that change the list of albums might have to refresh the medias !
* [X] Simplify the thunks that are only throwing a single action
* [X] Use a single dialog in the state (only one can be open!)
* [X] extract a selector for the albums !
* [X] Delete or create albums does not redirect (regression)
* [X] retire albumDatesUpdated in favour of using the albumAndMediasUpdated action (which needs to close the dialog)
* [ ] refactor the create dialog to match the new style