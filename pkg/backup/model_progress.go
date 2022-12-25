package backup

type ProgressEventType string

type ProgressEvent struct {
	Type      ProgressEventType // Type defines what's count, and size are about ; some might not be used.
	Count     int               // Count is the number of media
	Size      int               // Size is the sum of the size of the media concerned by this event
	Album     string            // Album is the folder name of the medias concerned by this event
	MediaType MediaType         // MediaType is the type of media ; only mandatory with 'uploaded' event
}
