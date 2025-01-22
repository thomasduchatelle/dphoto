export interface CreateAlbumRequest {
    name: string
    start: Date
    end: Date
    forcedFolderName: string
}

export class CreateAlbumHandler {
    constructor() {
    }

    handleCreateAlbum = (request: CreateAlbumRequest): Promise<void> => {
        return Promise.resolve()
    }
}