import {catalogReducerFunction, CurrentUserInsight, generateAlbumFilterOptions, initialCatalogState} from "./catalog-reducer";
import {Album, AlbumFilterEntry, CatalogViewerState, MediaType, MediaWithinADay, SharingType, UserDetails} from "./catalog-state";
import {CatalogViewerAction} from "./catalog-actions";

describe("CatalogViewerState", () => {
    const myselfUser: CurrentUserInsight = {picture: "my-face.jpg"};
    const herselfUser: UserDetails = {email: "her@self.com", name: "Herself", picture: "her-face.jpg"};
    const herselfOwner = "herself";

    const twoAlbums: Album[] = [
        {
            albumId: {owner: "myself", folderName: "jan-25"},
            name: "January 2025",
            start: new Date(2025, 0, 1),
            end: new Date(2025, 0, 31),
            totalCount: 42,
            temperature: 0.25,
            relativeTemperature: 1,
            sharedWith: [
                {
                    user: herselfUser,
                    role: SharingType.visitor,
                }
            ],
        },
        {
            albumId: {owner: herselfOwner, folderName: "feb-25"},
            name: "February 2025",
            start: new Date(2025, 1, 1),
            end: new Date(2025, 2, 0),
            totalCount: 12,
            temperature: 0.25,
            relativeTemperature: 1,
            ownedBy: {name: "Herself", users: [herselfUser]},
            sharedWith: [],
        },
    ]

    const twoAlbumsNoFilterOptions: AlbumFilterEntry = {
        criterion: {
            owners: []
        },
        avatars: [`${myselfUser.picture}`, `${herselfUser.picture}`],
        name: "All albums",
    };

    const someMedias: MediaWithinADay[] = [{
        day: new Date(2025, 0, 1), medias: [{
            id: "media-1",
            type: MediaType.IMAGE,
            time: new Date("2025-01-05T12:42:00Z"),
            uiRelativePath: "media-1/image.jpg",
            contentPath: "/media-1.jpg",
            source: "",
        }]
    }]

    const loadedStateWithTwoAlbums: CatalogViewerState = {
        allAlbums: twoAlbums,
        albumFilterOptions: [
            {
                criterion: {
                    selfOwned: true,
                    owners: [],
                },
                avatars: [myselfUser.picture ?? ""],
                name: "My albums",
            },
            twoAlbumsNoFilterOptions,
            {
                criterion: {
                    owners: [herselfOwner]
                },
                avatars: [herselfUser.picture ?? ""],
                name: herselfUser.name,
            },
        ],
        albumFilter: twoAlbumsNoFilterOptions,
        albums: twoAlbums,
        medias: someMedias,
        albumNotFound: false,
        mediasLoadedFromAlbumId: twoAlbums[0].albumId,
        albumsLoaded: true,
        mediasLoaded: true,
    };

    it("should add the loaded albums and medias to the state, and reset all status when receiving AlbumsAndMediasLoadedAction", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        const got = catalogReducer({
            ...initialCatalogState,
            albumNotFound: true,
            albumsLoaded: false,
            mediasLoaded: false,
        }, {
            type: "AlbumsAndMediasLoadedAction",
            albums: twoAlbums,
            medias: someMedias,
            selectedAlbum: twoAlbums[0],
        });

        expect(got).toEqual(loadedStateWithTwoAlbums)
    })

    it("should use 'All albums' filter even when it's the only selection available (only directly owned albums) when receiving AlbumsAndMediasLoadedAction", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        const got = catalogReducer(initialCatalogState, {
            type: "AlbumsAndMediasLoadedAction",
            albums: [twoAlbums[0]],
            medias: someMedias,
            selectedAlbum: twoAlbums[0],
        });

        const allAlbumFilter = {
            criterion: {
                owners: []
            },
            avatars: [myselfUser.picture ?? ""],
            name: "All albums",
        };
        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            albums: [twoAlbums[0]],
            allAlbums: [twoAlbums[0]],
            albumFilter: allAlbumFilter,
            albumFilterOptions: [allAlbumFilter],
        })
    })

    it("should show only directly owned album after the AlbumsFilteredAction", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        const got = catalogReducer(loadedStateWithTwoAlbums, {
            type: "AlbumsFilteredAction",
            criterion: {selfOwned: true, owners: []},
        });

        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[0],
            albums: [twoAlbums[0]],
        })
    })

    it("should show all albums when the filter moves back to 'All albums'", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        const got = catalogReducer({
            ...loadedStateWithTwoAlbums,
            albums: [],
        }, {
            type: "AlbumsFilteredAction",
            criterion: {owners: []},
        });

        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[1],
            albums: twoAlbums,
        })
    })

    it("should filter albums to those with a certain owner when the filter with that owner is selected", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        const got = catalogReducer({
            ...loadedStateWithTwoAlbums,
            albums: [],
        }, {
            type: "AlbumsFilteredAction",
            criterion: {owners: [herselfOwner]},
        });

        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[2],
            albums: [twoAlbums[1]],
        })
    })

    it("should only change the medias and loading status when reducing StartLoadingMediasAction, and clear errors", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        expect(catalogReducer({
            ...loadedStateWithTwoAlbums,
            albumNotFound: true,
            error: new Error("TEST previous error to clear"),
        }, {
            type: "StartLoadingMediasAction",
            albumId: twoAlbums[1].albumId,
        })).toEqual({
            ...loadedStateWithTwoAlbums,
            medias: [],
            loadingMediasFor: twoAlbums[1].albumId,
            mediasLoaded: false,
        })
    })

    it("should only change the medias and loading status when reducing MediasLoadedAction, and clear errors", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        expect(catalogReducer({
            ...loadedStateWithTwoAlbums,
            medias: [],
            mediasLoaded: false,
            loadingMediasFor: twoAlbums[1].albumId,
            albumNotFound: true,
            error: new Error("TEST previous error to clear"),
        }, {
            type: "MediasLoadedAction",
            albumId: twoAlbums[1].albumId,
            medias: someMedias,
        })).toEqual({
            ...loadedStateWithTwoAlbums,
            medias: someMedias,
            mediasLoadedFromAlbumId: twoAlbums[1].albumId,
        })
    })

    it("should ignore MediasLoadedAction if the medias are not for the expected album", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        expect(catalogReducer({
            ...loadedStateWithTwoAlbums,
            medias: [],
            loadingMediasFor: twoAlbums[0].albumId,
        }, {
            type: "MediasLoadedAction",
            albumId: twoAlbums[1].albumId,
            medias: someMedias,
        })).toEqual({
            ...loadedStateWithTwoAlbums,
            medias: [],
            loadingMediasFor: twoAlbums[0].albumId,
        })
    })

    it("should set the errors and clears medias and media loading status when reducing MediaFailedToLoadAction", () => {
        const testError = new Error("TEST loading error");

        const catalogReducer = catalogReducerFunction(myselfUser);
        expect(catalogReducer(initialCatalogState, {
            type: "MediaFailedToLoadAction",
            albums: twoAlbums,
            selectedAlbum: twoAlbums[0],
            error: testError,
        })).toEqual({
            ...loadedStateWithTwoAlbums,
            medias: [],
            mediasLoadedFromAlbumId: undefined,
            error: testError,
        })
    })

    it("should set the errors and clears medias and media loading status when reducing MediaFailedToLoadAction that hasn't the albums", () => {
        const testError = new Error("TEST loading error");

        const catalogReducer = catalogReducerFunction(myselfUser);
        expect(catalogReducer({
            ...initialCatalogState,
            allAlbums: twoAlbums,
        }, {
            type: "MediaFailedToLoadAction",
            selectedAlbum: twoAlbums[0],
            error: testError,
        })).toEqual({
            ...loadedStateWithTwoAlbums,
            medias: [],
            mediasLoadedFromAlbumId: undefined,
            error: testError,
        })
    })

    it("should update the list of albums and clear errors when AlbumsLoadedAction is received", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        expect(catalogReducer({
            ...loadedStateWithTwoAlbums,
            allAlbums: [twoAlbums[0]],
            albums: [twoAlbums[0]],
            error: new Error("TEST previous error to clear"),
            albumsLoaded: false,
        }, {
            type: "AlbumsLoadedAction",
            albums: twoAlbums,
        })).toEqual(loadedStateWithTwoAlbums)
    })

    it("should update the available filters and re-apply the selected filter when receiving AlbumsLoadedAction", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        expect(catalogReducer({
            ...loadedStateWithTwoAlbums,
            albumFilterOptions: [loadedStateWithTwoAlbums.albumFilterOptions[0]],
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[0],
            allAlbums: [twoAlbums[1]],
            albums: [],
            mediasLoadedFromAlbumId: twoAlbums[0].albumId, // no effect
        }, {
            type: "AlbumsLoadedAction",
            albums: twoAlbums,
            redirectTo: twoAlbums[0].albumId,
        })).toEqual({
            ...loadedStateWithTwoAlbums,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[0],
            allAlbums: twoAlbums,
            albums: [twoAlbums[0]],
        })
    })

    it("should remove the album filter if the redirectTo in AlbumsLoadedAction wouldn't be displayed", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        expect(catalogReducer({
            ...loadedStateWithTwoAlbums,
            allAlbums: [twoAlbums[0]],
            albums: [],
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[0],
        }, {
            type: "AlbumsLoadedAction",
            albums: twoAlbums,
            redirectTo: twoAlbums[1].albumId,
        })).toEqual({
            ...loadedStateWithTwoAlbums,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[1],
            allAlbums: twoAlbums,
            albums: twoAlbums,
        })
    })

    it("should open the sharing modal with the appropriate albumId and already-shared list", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        const action: CatalogViewerAction = {
            type: "OpenSharingModalAction",
            albumId: twoAlbums[0].albumId,
        };

        const expected: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                        role: SharingType.visitor,
                    }
                ],
            }
        };
        expect(catalogReducer(loadedStateWithTwoAlbums, action)).toEqual(expected);
    });

    it("should close the sharing modal by clearing the shareModel property", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        const action: CatalogViewerAction = {type: "CloseSharingModalAction"};

        const initial: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                        role: SharingType.visitor,
                    }
                ],
            }
        };
        expect(catalogReducer(initial, action)).toEqual(loadedStateWithTwoAlbums);
    });

    it("should add a new sharing entry and keep the modal open when receiving AddSharingAction", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        const initial: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                        role: SharingType.visitor,
                    }
                ],
            }
        };
        const newUser: UserDetails = {email: "bob@example.com", name: "Bob", picture: "bob-face.jpg"};
        const action: CatalogViewerAction = {
            type: "AddSharingAction",
            sharing: {
                user: newUser,
                role: SharingType.contributor,
            }
        };
        const expected: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: newUser,
                        role: SharingType.contributor,
                    },
                    {
                        user: herselfUser,
                        role: SharingType.visitor,
                    }
                ],
            }
        };
        expect(catalogReducer(initial, action)).toEqual(expected);
    });

    it("should replace an existing sharing entry for the same user when receiving AddSharingAction", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        const initial: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                        role: SharingType.visitor,
                    }
                ],
            }
        };
        // Add the same user with a different role: user is overridden and not added
        const action: CatalogViewerAction = {
            type: "AddSharingAction",
            sharing: {
                user: herselfUser,
                role: SharingType.contributor,
            }
        };
        const expected: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                        role: SharingType.contributor,
                    }
                ],
            }
        };
        expect(catalogReducer(initial, action)).toEqual(expected);
    });

    it("should remove a sharing entry by email and keep the modal open when receiving RemoveSharingAction", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        const bobEmail = "bob@example.com";
        const initial: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                        role: SharingType.visitor,
                    },
                    {
                        user: {email: bobEmail, name: "Bob", picture: "bob-face.jpg"},
                        role: SharingType.contributor,
                    }
                ],
            }
        };
        const action: CatalogViewerAction = {
            type: "RemoveSharingAction",
            email: bobEmail,
        };
        const expected: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                        role: SharingType.visitor,
                    }
                ],
            }
        };
        expect(catalogReducer(initial, action)).toEqual(expected);
    });

    it("should not change state when AddSharingAction is received and shareModal is closed", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        const initial: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            shareModal: undefined,
        };
        const newUser: UserDetails = {email: "bob@example.com", name: "Bob", picture: "bob-face.jpg"};
        const action: CatalogViewerAction = {
            type: "AddSharingAction",
            sharing: twoAlbums[0].sharedWith[0],
        };
        expect(catalogReducer(initial, action)).toEqual(initial);
    });

    it("should not change state when RemoveSharingAction is received and shareModal is undefined", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        const initial: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            shareModal: undefined,
        };
        const action: CatalogViewerAction = {
            type: "RemoveSharingAction",
            email: herselfUser.email,
        };
        expect(catalogReducer(initial, action)).toEqual(initial);
    });

    it("should not change state when RemoveSharingAction is received with an email not in sharedWith", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        const initial: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                        role: SharingType.visitor,
                    }
                ],
            }
        };
        const action: CatalogViewerAction = {
            type: "RemoveSharingAction",
            email: "notfound@example.com",
        };
        expect(catalogReducer(initial, action)).toEqual(initial);
    });

    it("should set the error field when receiving SharingModalErrorAction", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        const initial: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                        role: SharingType.visitor,
                    }
                ],
            }
        };
        const error = {type: "adding", message: "Failed to add user"} as const;
        const action: CatalogViewerAction = {
            type: "SharingModalErrorAction",
            error,
        };
        const expected: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                        role: SharingType.visitor,
                    }
                ],
                error,
            }
        };
        expect(catalogReducer(initial, action)).toEqual(expected);
    });

    it("should sort the sharedWith list alphabetically by user name when adding or opening sharings", () => {
        const catalogReducer = catalogReducerFunction(myselfUser);
        // Test AddSharingAction
        const initial: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                        role: SharingType.visitor,
                    }
                ],
            }
        };
        const tonyUser: UserDetails = {email: "tony@avenger.com", name: "", picture: "Tony-face.jpg"};
        const captainUser: UserDetails = {email: "captain@avenger.com", name: "Captain", picture: "captain-face.jpg"};

        const got = catalogReducer(catalogReducer(initial, {
            type: "AddSharingAction",
            sharing: {
                user: tonyUser,
                role: SharingType.visitor,
            }
        }), {
            type: "AddSharingAction",
            sharing: {
                user: captainUser,
                role: SharingType.visitor,
            }
        })
        const expected: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: captainUser,
                        role: SharingType.visitor,
                    },
                    {
                        user: herselfUser,
                        role: SharingType.visitor,
                    },
                    {
                        user: tonyUser,
                        role: SharingType.visitor,
                    }
                ],
            }
        };

        expect(got).toEqual(expected);
    });
})

describe("generateAlbumFilterOptions", () => {
    const currentUser = {picture: "my-face.jpg"};
    const herselfOwner = "herself";
    const herselfUser: UserDetails = {email: "her@self.com", name: "Herself", picture: "her-face.jpg"};
    const himselfOwner = "himself";
    const himself = {name: "Z-Himself", users: [{email: "him@self.com", name: "-", picture: "a-his-face.jpg"}]};
    const janBaseAlbum = {
        name: "January 2025",
        start: new Date(2025, 0, 1),
        end: new Date(2025, 1, 0),
        totalCount: 42,
        temperature: 0.25,
        relativeTemperature: 1,
        sharedWith: [],
    }
    const febBaseAlbum = {
        name: "February 2025",
        start: new Date(2025, 1, 1),
        end: new Date(2025, 2, 0),
        totalCount: 12,
        temperature: 0.25,
        relativeTemperature: 1,
        sharedWith: [],
    }

    it("should return the single option 'All albums' if no albums are available", () => {
        expect(generateAlbumFilterOptions(currentUser, [])).toEqual([{
            name: "All albums",
            criterion: {
                owners: []
            },
            avatars: [currentUser.picture],
        }])
    })

    it("should return the single option 'All albums' without picture if the current user doesn't have any", () => {
        expect(generateAlbumFilterOptions({}, [])).toEqual([{
            name: "All albums",
            criterion: {
                owners: []
            },
            avatars: [],
        }])
    })

    it("should return the single option 'All albums' if all albums are owned by the current user", () => {
        let janBaseAlbum = {
            name: "January 2025",
            start: new Date(2025, 0, 1),
            end: new Date(2025, 0, 31),
            totalCount: 42,
            temperature: 0.25,
            relativeTemperature: 1,
            sharedWith: [],
        };
        expect(generateAlbumFilterOptions(currentUser, [
            {
                ...janBaseAlbum,
                albumId: {owner: "myself", folderName: "jan-25"},
            },
        ])).toEqual([{
            name: "All albums",
            criterion: {
                owners: []
            },
            avatars: [currentUser.picture],
        }])
    })

    it("should return the single option 'All albums' if all albums are owned by another user", () => {
        expect(generateAlbumFilterOptions(currentUser, [
            {
                ...janBaseAlbum,
                albumId: {owner: "herself", folderName: "jan-25"},
                ownedBy: {name: "herself", users: [herselfUser]},
            },
        ])).toEqual([{
            name: "All albums",
            criterion: {
                owners: []
            },
            avatars: [currentUser.picture, herselfUser.picture],
        }])
    })

    it("should return the options 'My album', 'All albums', 'Her albums' if there is albums owned by current user and HER", () => {
        expect(generateAlbumFilterOptions(currentUser, [
            {
                ...janBaseAlbum,
                albumId: {owner: "myself", folderName: "jan-25"},
            },
            {
                ...febBaseAlbum,
                albumId: {owner: herselfOwner, folderName: "feb-25"},
                ownedBy: {name: herselfUser.name, users: [herselfUser]},
            },
        ])).toEqual([
            {
                name: "My albums",
                criterion: {
                    selfOwned: true,
                    owners: [],
                },
                avatars: [currentUser.picture],
            },
            {
                name: "All albums",
                criterion: {
                    owners: []
                },
                avatars: [currentUser.picture, herselfUser.picture],
            },
            {
                name: herselfUser.name,
                criterion: {
                    owners: [herselfOwner]
                },
                avatars: [herselfUser.picture],
            },
        ])
    })

    it("should return the options 'My album', 'All albums', 'Her albums' if there is albums owned by current user and HER even if HER doesn't have a picture", () => {
        expect(generateAlbumFilterOptions(currentUser, [
            {
                ...janBaseAlbum,
                albumId: {owner: "myself", folderName: "jan-25"},
            },
            {
                ...febBaseAlbum,
                albumId: {owner: "herself", folderName: "feb-25"},
                ownedBy: {name: herselfUser.name, users: [{...herselfUser, picture: undefined}]},
            },
        ])).toEqual([
            {
                name: "My albums",
                criterion: {
                    selfOwned: true,
                    owners: [],
                },
                avatars: [currentUser.picture],
            },
            {
                name: "All albums",
                criterion: {
                    owners: []
                },
                avatars: [currentUser.picture],
            },
            {
                name: herselfUser.name,
                criterion: {
                    owners: [herselfOwner]
                },
                avatars: [],
            },
        ])
    })

    it("should return the options 'All albums', 'Her albums', 'His albums' if all albums are owned by two other owners ; pictures are ordered by owner name", () => {
        expect(generateAlbumFilterOptions(currentUser, [
            {
                ...janBaseAlbum,
                albumId: {owner: himselfOwner, folderName: "jan-25"},
                ownedBy: himself,
            },
            {
                ...febBaseAlbum,
                albumId: {owner: herselfOwner, folderName: "jan-25"},
                ownedBy: {name: herselfUser.name, users: [herselfUser]},
            },
        ])).toEqual([
            {
                name: "All albums",
                criterion: {
                    owners: []
                },
                avatars: [currentUser.picture, herselfUser.picture, "a-his-face.jpg"],
            },
            {
                name: herselfUser.name,
                criterion: {
                    owners: [herselfOwner]
                },
                avatars: [herselfUser.picture],
            },
            {
                name: himself.name,
                criterion: {
                    owners: [himselfOwner]
                },
                avatars: [himself.users[0].picture],
            },
        ])
    })

    it("should return the options 'All albums', 'Her albums', 'His albums' in alphabetic order when albums comes in a different order", () => {
        expect(generateAlbumFilterOptions(currentUser, [
            {
                ...febBaseAlbum,
                albumId: {owner: herselfOwner, folderName: "jan-25"},
                ownedBy: {name: herselfUser.name, users: [herselfUser]},
            },
            {
                ...janBaseAlbum,
                albumId: {owner: himselfOwner, folderName: "jan-25"},
                ownedBy: himself,
            },
        ])).toEqual([
            {
                name: "All albums",
                criterion: {
                    owners: []
                },
                avatars: [currentUser.picture, herselfUser.picture, himself.users[0].picture],
            },
            {
                name: herselfUser.name,
                criterion: {
                    owners: [herselfOwner]
                },
                avatars: [herselfUser.picture],
            },
            {
                name: himself.name,
                criterion: {
                    owners: [himselfOwner]
                },
                avatars: [himself.users[0].picture],
            },
        ])
    })
})

export {}
