export const MediaPageMediasStateInit: MediaPageMediasState = {
  owner: "",
  folderName: "",
  medias: [],
}

export interface MediaPageMediasState {
  owner: string
  folderName: string
  medias: MediaRef[]
}

export interface MediaRef {
  encodedId: string
  filename: string
}