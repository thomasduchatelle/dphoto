# TODOs

1. implements the ports
2. integrate it with consumers (API and CLI) -> requires to add 'user' concept to the CLI
3. refactor to use a list of providers, each will fetch a list of albums that will be sorted after: OwnedAlbumsProvider, SharedAlbumsProvider,
   OfflineViewProviders, etc.
4. do the count of media as part of this view ; not the repository
5. start building a view that contains for a user:
   * owned albums + count + last media date + sharing
   * shared albums [by owner] + count + last media date
6. remove catalogaclview package
