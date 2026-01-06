---
applyTo: "web/**"
---

BEWARE - PROJECT `web/` AS BEEN DEPRECATED AND IS REPLACED BY `web-nextjs/`. Before making change to the `web/` project, make sure the intention was not to
change the new `web-nextjs/` project.

## Tree structure of `web/` project

The file structure is as follows:

* `components/`
    * `catalog-react/` - contains the React Components used to integrate the domain to the other components
    * `EditDateDialog` - ... as an example, the list of the components
* `core/catalog/` - "catalog" is the name of the domain
    * `<feature name in dash-case>/` - each feature has a folder containing related actions, thunks, and selectors
    * `language/` - ubiquitous language and definition of the State shared for the domain
    * `common/` - functionalities reused across most features
    * `actions.ts` - where the action interface and the partial reducer are registered
    * `thunks.ts` - where the thunks are registered
    * `adapters/api/` - where the REST API adapter are implemented

The same design principles of `web-nextjs/` applies to `web/` project.
