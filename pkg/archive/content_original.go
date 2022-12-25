package archive

import "time"

const DownloadUrlValidityDuration = 5 * time.Minute

// GetMediaOriginalURL returns a pre-signed URL to download the content. URL only valid a certain time.
func GetMediaOriginalURL(owner, mediaId string) (string, error) {
	key, err := repositoryPort.FindById(owner, mediaId)
	if err != nil {
		return "", err
	}

	return storePort.SignedURL(key, DownloadUrlValidityDuration)
}
