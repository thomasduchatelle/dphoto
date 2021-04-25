package backup

const (
	numberOfMediaType        = 3
	backupChannelsBufferSize = 32
)

var (
	ImageReaderThreadCount     = 4                      // ImageReaderThreadCount is the number of goroutines that will be used to read metadata from file contents
	UploadThreadCount          = 2                      // UploadThreadCount is the number of goroutines that will be used to backup batch of files
	DownloadThreadCount        = 2                      // DownloadThreadCount is the number of concurrent download from volume to local storage allowed
	UploadBatchSize            = 25                     // UploadBatchSize number of media to process as a batch (ideally, should be a multiple of 25: the dynamodb write batch size)
	LocalMediaPath             = "$HOME/.dphoto/medias" // LocalMediaPath is the path where medias are downloaded first before being read
	LocalBufferAreaSizeInOctet = 512 * 1024 * 1024      // LocalBufferAreaSizeInOctet is the size, in octet, dphoto can use to download photos
	OnlineBackupLocation       string                   // OnlineBackupLocation is the S£ bucket name where document must be uplaoded.
)
