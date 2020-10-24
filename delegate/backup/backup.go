package backup

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"path"
	"sync"
	"time"
)

const (
	numberOfMediaType        = 3
	backupChannelsBufferSize = 32
)

type Runner struct {
	context           *log.Entry
	BackupId          string
	VolumeId          string
	temporaryLocalDir string

	counter Counter

	foundChannel          chan FoundMedia
	readyForBackupChannel chan LocalMedia
	CompletionChannel     chan Report

	errorsMutex *sync.Mutex
	errors      []error
}

func StartBackupRunner(volume RemovableVolume) (*Runner, error) {
	if len(volume.MountPaths) == 0 {
		return nil, errors.Errorf("volume must have at least 1 mount point.")
	}

	withContext := log.WithFields(log.Fields{
		"VolumeId":   volume.UniqueId,
		"MountPaths": volume.MountPaths,
	})

	snapshot, err := VolumeRepository.RestoreLastSnapshot(volume.UniqueId)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to restore previous snapshot.")
	}
	lastVolumeSnapshot := make(map[string]int)
	for _, m := range snapshot {
		lastVolumeSnapshot[m.RelativePath] = m.Size
	}

	backupId := fmt.Sprintf("%s_%s", volume.UniqueId, time.Now().Format("20060102_150405"))

	temporaryDir, err := temporaryMediaPath(backupId)
	if err != nil {
		return nil, err
	}

	runner := Runner{
		context:           withContext.WithField("BackupId", backupId),
		BackupId:          backupId,
		VolumeId:          volume.UniqueId,
		temporaryLocalDir: temporaryDir,
		counter:           Counter{},

		foundChannel:          make(chan FoundMedia, backupChannelsBufferSize),
		readyForBackupChannel: make(chan LocalMedia, backupChannelsBufferSize),
		CompletionChannel:     make(chan Report, 1),

		errorsMutex: &sync.Mutex{},
	}

	runner.start(volume.MountPaths[0], lastVolumeSnapshot)

	return &runner, nil
}

func (r *Runner) start(mountPath string, lastVolumeSnapshot map[string]int) {
	analiseWorkerGroup := sync.WaitGroup{}
	var signatures [ImageReaderThread][]SimpleMediaSignature

	go func() {
		err := FileHandler.FindMediaRecursively(mountPath, r.foundChannel)
		if err != nil {
			r.appendError(err)
		}
		close(r.foundChannel)
	}()

	for thread := 0; thread < ImageReaderThread; thread++ {
		runnerThread := thread
		analiseWorkerGroup.Add(1)
		go func() {
			defer analiseWorkerGroup.Done()
			workerSignatures, err := r.copyToLocalAndAnalyse(lastVolumeSnapshot, r.foundChannel, r.readyForBackupChannel)
			if err != nil {
				r.appendError(err)
			} else {
				signatures[runnerThread] = workerSignatures
			}
		}()
	}

	go func() {
		analiseWorkerGroup.Wait()
		close(r.readyForBackupChannel)
	}()

	go func() {
		err := OnlineStorage.BackupOnline(r.readyForBackupChannel)
		if err != nil {
			r.appendError(err)
		}

		r.completionRoutine(&analiseWorkerGroup, &signatures)
	}()
}

func (r *Runner) copyToLocalAndAnalyse(ignoreList map[string]int, origin chan FoundMedia, dest chan LocalMedia) (signature []SimpleMediaSignature, err error) {
	for file := range origin {
		signature = append(signature, file.SimpleSignature)

		if err != nil {
			// continue to consume all messages in the event of an error (prevent dead-lock)
			continue
		}

		if size, present := ignoreList[file.SimpleSignature.RelativePath]; !present || file.SimpleSignature.Size != size {
			r.counter.incrementFoundCounter(file.Type)
			imagePath := path.Join(r.temporaryLocalDir, file.SimpleSignature.RelativePath)

			var mediaHash string
			mediaHash, err = FileHandler.CopyToLocal(file.LocalAbsolutePath, imagePath)
			if err != nil {
				r.context.WithError(err).Errorln(err)
				continue
			}

			var details *MediaDetails
			if file.Type == IMAGE {
				var detailsError error
				details, detailsError = ImageDetailsReader.ReadImageDetails(imagePath)
				if detailsError != nil {
					r.context.WithField("Image", file.LocalAbsolutePath).WithError(err).Warn("Failed reading image details")
				}
			}

			dest <- LocalMedia{
				Type:              file.Type,
				LocalAbsolutePath: imagePath,
				Signature: FullMediaSignature{
					Sha256: mediaHash,
					Size:   0,
				},
				Details: details,
			}
		}
	}

	return
}

func (r *Runner) completionRoutine(analiseWorkerGroup *sync.WaitGroup, signatures *[ImageReaderThread][]SimpleMediaSignature) {
	success := len(r.errors) == 0
	if success {
		r.context.Infof("%d media backed-up: %d images, %d videos", r.counter.GetFoundCount(), r.counter.GetFound(IMAGE), r.counter.GetFound(VIDEO))

		numberOfSignatures := 0
		for _, s := range signatures {
			numberOfSignatures += len(s)
		}

		allSignatures := make([]SimpleMediaSignature, 0, numberOfSignatures)
		for _, signatureSlice := range signatures {
			allSignatures = append(allSignatures, signatureSlice...)
		}

		err := VolumeRepository.StoreSnapshot(r.VolumeId, r.BackupId, allSignatures)
		if err != nil {
			r.context.WithError(err).Error("Failed to store volume snapshot, next backup will go through all media again.")
		}
	}

	r.CompletionChannel <- Report{
		Success: success,
		Errors:  r.errors,
		Counter: &r.counter,
	}
}

func (r *Runner) appendError(err error) {
	r.errorsMutex.Lock()
	defer r.errorsMutex.Unlock()

	r.errors = append(r.errors, err)
}
