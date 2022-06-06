import {MustBeAuthenticated} from "../../../../core/application";
import {MediaPageMediasState} from "./domain";

interface RestMedia {
  id: string
  filename: string
  time: string
}

export class MediaPageLogic {
  constructor(
    private mustBeAuthenticated: MustBeAuthenticated,
    private setState: (value: (MediaPageMediasState)) => void,
  ) {
  }

  public findMediasWithCache = (owner: string, folderName: string): Promise<void> => {
    return this.mustBeAuthenticated.authenticatedAxios.get<RestMedia[]>(`/api/v1/owners/${owner}/albums/${folderName}/medias`)
      .then(resp => {
        this.setState({
          owner,
          folderName,
          medias: resp.data
            .sort((a, b) => new Date(b.time).getTime() - new Date(a.time).getTime())
            .map(media => ({
              encodedId: media.id,
              filename: media.filename,
            })),
        })
      })
      .then()
  }

  // getModel is a synchronous function that will give an answer even if it triggers a reloading later on
  public getModel = (owner: string | undefined, album: string | undefined, encodedId: string | undefined, filename: string | undefined, mediasState: MediaPageMediasState): {
    backToAlbumLink: string
    imgSrc: string
    previousMediaLink: string
    nextMediaLink: string
  } => {
    let previousMediaLink = "", nextMediaLink = ""

    if (mediasState.owner === owner && mediasState.folderName === album && mediasState.medias) {
      let index = mediasState.medias.findIndex(media => media.encodedId === encodedId)
      if (index < 0) {
        index = 0
      }

      if (index > 0) {
        previousMediaLink = `/albums/${owner}/${album}/${mediasState.medias[index - 1].encodedId}/${mediasState.medias[index - 1].filename}`
      }

      if (index + 1 < mediasState.medias.length) {
        nextMediaLink = `/albums/${owner}/${album}/${mediasState.medias[index + 1].encodedId}/${mediasState.medias[index + 1].filename}`
      }
    }

    return {
      backToAlbumLink: `/albums/${owner}/${album}`,
      imgSrc: `/api/v1/owners/${owner}/medias/${encodedId}/${filename}`,
      previousMediaLink: previousMediaLink,
      nextMediaLink: nextMediaLink,
    }
  }
}