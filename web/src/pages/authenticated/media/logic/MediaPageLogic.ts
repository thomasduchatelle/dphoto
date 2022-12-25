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
    currentIsImage: boolean
    previousMediaLink: string
    previousMediaSrc: string
    previousIsImage: boolean
    nextMediaLink: string
    nextMediaSrc: string
    nextIsImage: boolean
  } => {
    let previousMediaLink = "", nextMediaLink = "", previousMediaSrc = "", nextMediaSrc = "", previousIsImage = false,
      nextIsImage = false

    if (mediasState.owner === owner && mediasState.folderName === album && mediasState.medias) {
      let index = mediasState.medias.findIndex(media => media.encodedId === encodedId)
      if (index < 0) {
        index = 0
      }

      if (index > 0) {
        const media = mediasState.medias[index - 1];
        previousMediaLink = `/albums/${owner}/${album}/${media.encodedId}/${media.filename}`
        previousMediaSrc = `/api/v1/owners/${owner}/medias/${media.encodedId}/${media.filename}?access_token=${this.mustBeAuthenticated.accessToken}`
        previousIsImage = this.isAnImage(media.filename)
      }

      if (index + 1 < mediasState.medias.length) {
        const media = mediasState.medias[index + 1];
        nextMediaLink = `/albums/${owner}/${album}/${media.encodedId}/${media.filename}`
        nextMediaSrc = `/api/v1/owners/${owner}/medias/${media.encodedId}/${media.filename}?access_token=${this.mustBeAuthenticated.accessToken}`
        nextIsImage = this.isAnImage(media.filename)
      }
    }

    return {
      backToAlbumLink: `/albums/${owner}/${album}`,
      imgSrc: `/api/v1/owners/${owner}/medias/${encodedId}/${filename}?access_token=${this.mustBeAuthenticated.accessToken}`,
      currentIsImage: this.isAnImage(filename),
      previousMediaLink,
      previousMediaSrc,
      previousIsImage,
      nextMediaLink,
      nextMediaSrc,
      nextIsImage,
    }
  }

  private isAnImage(filename: string | undefined): boolean {
    return [".jpg", ".jpeg", ".png"].filter(ext => filename?.toLowerCase().endsWith(ext)).length > 0;
  }
}